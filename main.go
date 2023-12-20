package main

import (
	"fmt"

	"github.com/SleepWalker/sinx/iface"
	"github.com/SleepWalker/sinx/snet"
)

func main() {
	server := snet.NewServer()
	//添加自定义router
	server.AddRouter(&PingRouter{})
	server.Serve()
}

// ping 自定义路由
type PingRouter struct {
	snet.BaseRouter
}

func (pr *PingRouter) Handle(request iface.IRequest) {
	fmt.Println("Call Router Handle...")
	// 1.读取客户端数据
	fmt.Println("recv from client: msgID = ", request.GetMsgID(), ", data = ", string(request.GetData()))

	// 2.回显ping
	err := request.GetConnection().SendMsg(request.GetMsgID(), []byte("ping...ping...\n"))
	if err != nil {
		fmt.Println("cal back handle error: ", err)
	}
}
