package snet

import (
	"fmt"
	"net"

	"github.com/SleepWalker/sinx/iface"
)

type Server struct {
	Name      string `info:"服务器名称"`
	IPVersion string `info:"IP版本"`
	IP        string `info:"服务器监听的IP"`
	Port      int    `info:"服务器监听的端口"`
}

func (s *Server) Start() {
	// 1.获取TCP的Addr
	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("net resolve addr error: ", err)
	}

	// 2.监听地址
	listener, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		fmt.Println("listen TCP error: ", err)
		return
	}

	fmt.Println("start Sinx server success, ", s.Name, " is listening...")

	// 3.阻塞，等待客户端连接，处理客户端业务
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("accept error: ", err)
			continue
		}

		// 做一个回显业务，最大512字节
		go func() {
			buf := make([]byte, 512)
			for {
				len, err := conn.Read(buf)
				if err != nil {
					fmt.Println("read buffer error: ", err)
					continue
				}

				// 回显，发送给客户端
				if _, err := conn.Write(buf[:len]); err != nil {
					fmt.Println("write buffer error: ", err)
					continue
				}
			}
		}()
	}

}

func (s *Server) Stop() {
	// TODO 资源、状态、链接信息 停止或回收

}

func (s *Server) Serve() {
	//Serve要处理其他业务，不能再Start中阻塞，故开启goroutine
	go s.Start()

	// TODO 启动服务器后的额外业务

	//阻塞
	select {}
}

/**
* 返回一个Server对象
**/
func NewServer(name string) iface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8888,
	}

	return s
}
