package router

import (
	"SleepXLink/iface"
	"SleepXLink/network"
	"fmt"
)

/****************************************
@Author : SleepWalker
@Description:
@Time : 2024/1/2 17:01
****************************************/

// hello 自定义路由
type HelloRouter struct {
	network.BaseRouter
}

func (hr *HelloRouter) Handle(request iface.IRequest) {
	fmt.Println("Call Hello Router Handle...")
	// 1.读取客户端数据
	fmt.Println("receive from client: msgID = ", request.GetMsgID(), ", data = ", string(request.GetData()))

	// 2.回显ping
	err := request.GetConnection().SendMsg(500, []byte("Hello SleepXLink!\n"))
	if err != nil {
		fmt.Println("cal back handle error: ", err)
	}
}
