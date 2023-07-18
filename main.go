package main

import (
	"fmt"
	"github.com/invictus555/auto_codes/greeting_service_v1/kitex_gen/greeting"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/server"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/invictus555/auto_codes/greeting_service_v1/kitex_gen/greeting/greetingservice"
	consul "github.com/kitex-contrib/registry-consul"
)

const (
	consulAddr  = "127.0.0.1:8500"
	serviceName = "greeting.service"
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

	servicePort := genRandomPort(8800, 8888)
	fmt.Println("greeting.service port :", servicePort)
	svr := greetingservice.NewServer(
		NewGreetingServiceImpl(servicePort),
		server.WithRegistry(r),
		server.WithServiceAddr(&net.TCPAddr{
			Port: servicePort, // 可通过修改端口模拟多实例服务
			IP:   net.IPv4zero,
		}),
		server.WithRegistryInfo(&registry.Info{ //设置注册信息
			Weight:      3,           // 权重,若测试带权类负载均衡算法,不同实例的权重应该不一样
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

// 获取随机端口,方便同一测试机多实例部署
func genRandomPort(min, max int) int {
	if max < min {
		min, max = max, min
	}

	rand.Seed(time.Now().Unix())
	port := int(rand.Int31n(int32(max-min))) + min
	return port
}
