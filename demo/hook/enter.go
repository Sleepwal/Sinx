package hook

import (
	"SleepXLink/iface"
	"fmt"
)

/****************************************
@Author : SleepWalker
@Description:
@Time : 2024/1/2 17:02
****************************************/

// 创建连接的时候执行
func DoConnectionBegin(conn iface.IConnection) {
	fmt.Println("DoConnectionBegin is Called ... ")

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

	fmt.Println("DoConnectionLost is Called ... ")
}
