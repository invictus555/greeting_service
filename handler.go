package main

import (
	"context"
	greeting "github.com/invictus555/auto_codes/greeting_service_v1/kitex_gen/greeting"
	"strconv"
)

// GreetingServiceImpl implements the last service interface defined in the IDL.
type GreetingServiceImpl struct {
	servicePort int
}

// SayHello implements the GreetingServiceImpl interface.
func (s *GreetingServiceImpl) SayHello(ctx context.Context, req *greeting.Request) (resp *greeting.Response, err error) {
	// TODO: Your code here...
	err = nil
	resp = &greeting.Response{
		Message: req.Message + ", service comes from port = " + strconv.Itoa(s.servicePort),
	}
	return
}
