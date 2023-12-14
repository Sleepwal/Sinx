package snet

import (
	"fmt"
	"net"

	"github.com/SleepWalker/sinx/iface"
)

type Connection struct {
	// 当前链接的socket
	Conn *net.TCPConn

	// 链接ID
	ConnID uint32

	// 当前的链接状态
	isClosed bool

	// 当前链接绑定的处理业务方法
	handleAPI iface.HandleFunc

	// 通道: 告知当前链接已经停止
	ExitChan chan bool
}

// 初始化
func NewConnection(conn *net.TCPConn, connID uint32, callbackAPI iface.HandleFunc) *Connection {
	c := &Connection{
		Conn:      conn,
		ConnID:    connID,
		handleAPI: callbackAPI,
		isClosed:  false,
		ExitChan:  make(chan bool, 1),
	}

	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Connection StartReader is running...")
	defer fmt.Println("Connection ID = ", c.ConnID, " Reader is exit, address is ", c.RemoteAddr().String())
	defer c.Stop()

	buf := make([]byte, 512)
	for {
		//读取客户端数据到buf
		len, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("Connection Read error: ", err)
			continue
		}

		//调用绑定的handleAPI
		if err := c.handleAPI(c.Conn, buf, len); err != nil {
			fmt.Println("Connection handleAPI error: ", err)
			break
		}
	}
}

// 启动链接
func (c *Connection) Start() {
	fmt.Println("Connection Start... ConnID = ", c.ConnID)
	// 从当前链接读数据
	go c.StartReader()

	//TODO 写数据
}

// 停止链接
func (c *Connection) Stop() {
	fmt.Println("Connection Stop... ConnID = ", c.ConnID)

	// 已关闭
	if c.isClosed {
		return
	}

	c.isClosed = true
	c.Conn.Close()    // 关闭socket连接
	close(c.ExitChan) // 回收资源
}

// 获取当前链接绑定的socket
func (c *Connection) GetTCPConnect() *net.TCPConn {
	return c.Conn
}

// 获取当前链接模块的链接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取远程客户端的TCP状态: IP、Port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送数据给远程的客户端
func (c *Connection) Send(data []byte) error {
	return nil
}