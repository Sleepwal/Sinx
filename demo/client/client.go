package main

import (
	"fmt"
	"net"
	"time"

	"SleepXLink/network"
)

/*
模拟客户端
*/
func main() {
	fmt.Println("Client Test ... start")
	//3秒之后发起测试请求，给服务端开启服务的机会
	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	for i := 0; i < 5; i++ {
		//发封包message消息
		msg, _ := network.Pack(network.NewMessage(uint32(i)%2, []byte(fmt.Sprintf("SleepXLink Client%d Test Message", i))))
		_, err := conn.Write(msg)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		receiveMsg, err := network.UnPack(conn)
		if err != nil {
			fmt.Println("unpack error err ", err)
			return
		}
		fmt.Println("==> Recv Msg: ID=", receiveMsg.GetMsgId(),
			", len=", receiveMsg.GetDataLen(),
			", data=", string(receiveMsg.GetData()))

		time.Sleep(1 * time.Second)
	}
}
