package snet

import (
	"errors"
	"fmt"
	"io"
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
	// 通道: 告知当前链接已经停止
	ExitChan chan bool
	// 消息管理模块
	MsgHandle iface.IMsgHandle
}

// 初始化
func NewConnection(conn *net.TCPConn, connID uint32, handler iface.IMsgHandle) *Connection {
	c := &Connection{
		Conn:      conn,
		ConnID:    connID,
		MsgHandle: handler,
		isClosed:  false,
		ExitChan:  make(chan bool, 1),
	}

	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Connection StartReader is running...")
	defer fmt.Println("Connection ID = ", c.ConnID, " Reader is exit, address is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		// 读取客户端消息的Head，8个字节
		headBuf := make([]byte, GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnect(), headBuf); err != nil {
			fmt.Println("Connection Read head error: ", err)
			break
		}

		// 拆包，得到 MsgID 和 MsgData
		msg, err := UnPack(headBuf)
		if err != nil {
			fmt.Println("Message UnPack error: ", err)
			break
		}

		// 根据DataLen，接着读取消息中的Data
		if msg.GetDataLen() > 0 {
			dataBuf := make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnect(), dataBuf); err != nil {
				fmt.Println("Connection Read Message data error: ", err)
				break
			}

			msg.SetData(dataBuf)
		}

		// 得到当前conn数据的request请求数据
		req := &Requset{
			conn: c,
			msg:  msg,
		}

		// 执行注册的路由方法
		go c.MsgHandle.DoMsgHandler(req)
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

/**
* 给客户端发送数据，先封包，再发送
**/
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection is closed")
	}

	// 封包 len|id|data
	binaryMsg, err := Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Println("Message Pack error: ", err)
		return errors.New("msg Pack error")
	}

	// 发送给客户端
	if _, err := c.Conn.Write(binaryMsg); err != nil {
		fmt.Println("Connection Write Msg error: ", err)
		return errors.New("Connection Write Msg error")
	}

	return nil
}
