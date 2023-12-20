package iface

import "net"

// 链接模块的抽象层
type IConnection interface {
	// 启动链接
	Start()

	// 停止链接
	Stop()

	// 获取当前链接绑定的socket
	GetTCPConnect() *net.TCPConn

	// 获取当前链接模块的链接ID
	GetConnID() uint32

	// 获取远程客户端的TCP状态: IP、Port
	RemoteAddr() net.Addr

	// 发送数据给远程的客户端
	SendMsg(msgId uint32, data []byte) error
}

// 处理链接业务
type HandleFunc func(*net.TCPConn, []byte, int) error
