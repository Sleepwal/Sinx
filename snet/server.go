package snet

import (
	"fmt"
	"net"

	"github.com/SleepWalker/sinx/iface"
	"github.com/SleepWalker/sinx/utils"
)

type Server struct {
	Name      string `info:"服务器名称"`
	IPVersion string `info:"IP版本"`
	IP        string `info:"服务器监听的IP"`
	Port      int    `info:"服务器监听的端口"`
	// server注册的链接对应的处理业务
	Router iface.IRouter `info:"路由"`
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

	var cid uint32 = 0

	// 3.阻塞，等待客户端连接，处理客户端业务
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("accept error: ", err)
			continue
		}

		// 处理新连接，用链接模块处理
		handleConn := NewConnection(conn, cid, s.Router)
		cid++

		// 启动一个goroutine处理业务
		go handleConn.Start()
	}

}

func (s *Server) Stop() {
	// TODO 资源、状态、链接信息 停止或回收

}

func (s *Server) Serve() {
	fmt.Println("[Sinx]---Serve start!")
	fmt.Println("[Sinx]---Server Name:", utils.GlobalObject.Name)
	fmt.Println("[Sinx]---Server IP:", utils.GlobalObject.Host)
	fmt.Println("[Sinx]---Server Port:", utils.GlobalObject.Port)
	fmt.Println("[Sinx]---Server Version:", utils.GlobalObject.Version,
		", Server MaxConn:", utils.GlobalObject.MaxConn,
		", Server MaxPackageSize:", utils.GlobalObject.MaxPackageSize)

	//Serve要处理其他业务，不能再Start中阻塞，故开启goroutine
	go s.Start()

	// TODO 启动服务器后的额外业务

	//阻塞
	select {}
}

func (s *Server) AddRouter(router iface.IRouter) {
	s.Router = router
	fmt.Println("[Server]---AddRouter success!	")
}

/**
* 返回一个Server对象
**/
func NewServer(name string) iface.IServer {
	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.Port,
		Router:    nil,
	}

	return s
}
