package util

import (
	"math/rand"
	"net"
	"time"
)

// GetAvailablePort 获取随机可用端口,方便同一测试机多实例部署
func GetAvailablePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}

	ret, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}

	defer ret.Close()
	return ret.Addr().(*net.TCPAddr).Port, nil
}

// GetRandomWeight 获取随机权重, 方便观察带权负载均衡算法
func GetRandomWeight(min, max int) int {
	if max < min {
		min, max = max, min
	}

	rand.Seed(time.Now().Unix())
	port := int(rand.Int31n(int32(max-min))) + min
	return port
}
