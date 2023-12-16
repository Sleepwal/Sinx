package main

import (
	"fmt"

	"github.com/SleepWalker/sinx/iface"
	"github.com/SleepWalker/sinx/snet"
)

func main() {
	server := snet.NewServer("V0.3")
	//添加自定义router
	server.AddRouter(&PingRouter{})
	server.Serve()
}

// ping 自定义路由
type PingRouter struct {
	snet.BaseRouter
}

func (pr *PingRouter) PreHandle(request iface.IRequest) {
	fmt.Println("Call Router PreHandle...")
	_, err := request.GetConnection().GetTCPConnect().Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("cal back pre handle error: ", err)
	}
}

func (pr *PingRouter) Handle(request iface.IRequest) {
	fmt.Println("Call Router Handle...")
	_, err := request.GetConnection().GetTCPConnect().Write([]byte("ping ping...\n"))
	if err != nil {
		fmt.Println("cal back handle error: ", err)
	}
}

func (pr *PingRouter) PostHandle(request iface.IRequest) {
	fmt.Println("Call Router PostHandle...")
	_, err := request.GetConnection().GetTCPConnect().Write([]byte("after ping...\n"))
	if err != nil {
		fmt.Println("cal back post handle error: ", err)
	}
}
