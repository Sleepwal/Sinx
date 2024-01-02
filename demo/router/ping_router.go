package router

import (
	"SleepXLink/iface"
	"SleepXLink/network"
	"fmt"
)

/****************************************
@Author : SleepWalker
@Description:
@Time : 2024/1/2 17:00
****************************************/

// ping 自定义路由
type PingRouter struct {
	network.BaseRouter
}

func (pr *PingRouter) Handle(request iface.IRequest) {
	fmt.Println("Call Ping Router Handle...")
	// 1.读取客户端数据
	fmt.Println("recv from client: msgID = ", request.GetMsgID(), ", data = ", string(request.GetData()))

	// 2.回显ping
	err := request.GetConnection().SendMsg(200, []byte("ping...ping...\n"))
	if err != nil {
		fmt.Println("cal back handle error: ", err)
	}
}
