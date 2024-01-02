package net

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func TestDataPack(t *testing.T) {
	// === 模拟服务器 ===
	listenner, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	defer listenner.Close()

	go func() {
		for {
			// 读取数据
			conn, err := listenner.Accept()
			if err != nil {
				fmt.Println("accept failed, err:", err)
				return
			}

			// 处理客户端请求
			go func(conn net.Conn) {
				for {
					msg, err := UnPack(conn)
					if err != nil {
						fmt.Println("unpack failed, err:", err)
						return
					}

					fmt.Println("----> receive msg:", msg.GetMsgId(), msg.GetDataLen(), ", data: ", string(msg.GetData()))
				}
			}(conn)
		}
	}()

	// === 模拟客户端 ===
	go func() {
		conn, err := net.Dial("tcp", "127.0.0.1:8888")
		if err != nil {
			fmt.Println("dial failed, err:", err)
			return
		}
		defer conn.Close()

		data := []byte("Hello World")
		package1, err := Pack(&Message{
			ID:      1,
			DataLen: uint32(len(data)),
			Data:    data,
		})
		if err != nil {
			fmt.Println("pack package1 failed, err:", err)
			return
		}

		data2 := []byte("Hello Go")
		package2, err := Pack(&Message{
			ID:      2,
			DataLen: uint32(len(data2)),
			Data:    data2,
		})
		if err != nil {
			fmt.Println("pack package2 failed, err:", err)
			return
		}

		data3 := []byte("Golang")
		package3, err := Pack(&Message{
			ID:      3,
			DataLen: uint32(len(data3)),
			Data:    data3,
		})
		if err != nil {
			fmt.Println("pack package2 failed, err:", err)
			return
		}

		// 粘包发送
		sendData := append(append(package1, package2...), package3...)
		conn.Write(sendData)
	}()

	time.Sleep(2 * time.Second)
}
