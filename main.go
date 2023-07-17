package main

import (
	"github.com/invictus555/auto_codes/greeting_service_v1/kitex_gen/greeting/greetingservice"
	"log"
)

func main() {
	svr := greetingservice.NewServer(new(GreetingServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
