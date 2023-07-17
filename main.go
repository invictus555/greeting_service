package main

import (
	"log"
	"net"

	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/server"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/invictus555/auto_codes/greeting_service_v1/kitex_gen/greeting/greetingservice"
	consul "github.com/kitex-contrib/registry-consul"
)

const (
	Port        = 8088
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

	svr := greetingservice.NewServer(
		new(GreetingServiceImpl),
		server.WithRegistry(r),
		server.WithServiceAddr(&net.TCPAddr{
			IP: net.IPv4(192, 168, 31, 76), Port: Port, // 为了在一台机器上多实例测试，可以更改端口
		}),
		server.WithRegistryInfo(&registry.Info{ //设置注册信息
			Weight:      2,           // 权重,若为测试带权算法,此处可以修改成不同的值
			ServiceName: serviceName, // 服务名称
		}),
	)

	if err := svr.Run(); err != nil { // 启动服务器
		log.Println(err.Error())
	}
}
