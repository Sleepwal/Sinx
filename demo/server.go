package main

import (
	"SleepXLink/network"
	"xlink-demo/hook"
	"xlink-demo/router"
)

func main() {
	server := network.NewServer()

	//注册链接hook回调函数
	server.SetOnConnStart(hook.DoConnectionBegin)
	server.SetOnConnStop(hook.DoConnectionLast)

	//添加自定义router
	server.AddRouter(0, &router.HelloRouter{})
	server.AddRouter(1, &router.HelloRouter{})
	server.Serve()
}
