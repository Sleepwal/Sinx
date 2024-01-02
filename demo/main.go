package main

import (
	"fmt"

	"SleepXLink/iface"
	"SleepXLink/network"
)

func main() {
	server := network.NewServer()

	//注册链接hook回调函数
	server.SetOnConnStart(DoConnectionBegin)
	server.SetOnConnStop(DoConnectionLast)

	//添加自定义router
	server.AddRouter(0, &PingRouter{})
	server.AddRouter(1, &HelloRouter{})
	server.Serve()
}

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

// 创建连接的时候执行
func DoConnectionBegin(conn iface.IConnection) {
	fmt.Println("DoConnecionBegin is Called ... ")

	//=============设置两个链接属性，在连接创建之后===========
	fmt.Println("Set conn Name, Home done!")
	conn.SetProperty("Name", "SleepWalker")
	conn.SetProperty("Home", "https://github.com/Sleepwal")
	//===================================================

	err := conn.SendMsg(2, []byte("DoConnection BEGIN..."))
	if err != nil {
		fmt.Println(err)
	}
}

// 连接断开的时候执行
func DoConnectionLast(conn iface.IConnection) {
	//============在连接销毁之前，查询conn的Name，Home属性=====
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Conn Property Name = ", name)
	}

	if home, err := conn.GetProperty("Home"); err == nil {
		fmt.Println("Conn Property Home = ", home)
	}
	//===================================================

	fmt.Println("DoConneciotnLost is Called ... ")
}
