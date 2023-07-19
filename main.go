package main

import (
	"fmt"
	"github-Projs/greeting_service/util"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/server"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/invictus555/auto_codes/greeting_service_v1/kitex_gen/greeting"
	"github.com/invictus555/auto_codes/greeting_service_v1/kitex_gen/greeting/greetingservice"
	consul "github.com/kitex-contrib/registry-consul"
	"log"
	"net"
)

const (
	consulAddr  = "127.0.0.1:8500"
	serviceName = "hello.greeting.service"
)

func main() {
	r, err := consul.NewConsulRegister(
		consulAddr,
		consul.WithCheck(&consulApi.AgentServiceCheck{ // 设置健康检查
			Interval:                       "5s",
			Timeout:                        "3s",
			DeregisterCriticalServiceAfter: "1m",
		}))
	if err != nil {
		panic("register consul failed")
	}

	weight := util.GetRandomWeight(1, 5)
	servicePort, err := util.GetAvailablePort()
	if err != nil {
		panic("no available port, " + err.Error())
	}

	fmt.Printf("<hello.greeting.service> listen at port = %d, weight = %d\n", servicePort, weight)
	svr := greetingservice.NewServer(
		NewGreetingServiceImpl(servicePort),
		server.WithRegistry(r),
		server.WithServiceAddr(&net.TCPAddr{
			Port: servicePort, // 可通过修改端口模拟多实例服务
			IP:   net.IPv4zero,
		}),
		server.WithRegistryInfo(&registry.Info{ //设置注册信息
			Weight:      weight,      // 权重,若测试带权类负载均衡算法,不同实例的权重应该不一样
			ServiceName: serviceName, // 服务名称
		}),
	)

	if err := svr.Run(); err != nil { // 启动服务器
		log.Println(err.Error())
	}
}

func NewGreetingServiceImpl(port int) greeting.GreetingService {
	impl := new(GreetingServiceImpl)
	impl.servicePort = port
	return impl
}
