package net

import (
	"errors"
	"fmt"
	"net"
	"sync"

	"SleepXLink/iface"
	"SleepXLink/utils"
)

type Connection struct {
	// 当前Conn属于的Server
	TcpServer iface.IServer
	// 当前链接的socket
	Conn *net.TCPConn
	// 链接ID
	ConnID uint32
	// 当前的链接状态
	isClosed bool
	// 通道: 告知当前链接已经停止（Reader给Writer）
	ExitChan chan bool
	// 无缓冲管道，用于读、写Goroutine之间的通信
	msgChan chan []byte
	//有缓冲管道，用于读、写两个goroutine之间的消息通信
	msgBuffChan chan []byte
	// 消息管理模块
	MsgHandle iface.IMsgHandle

	//链接属性
	property map[string]interface{}
	//保护链接属性修改的锁
	propertyLock sync.RWMutex
}

// 初始化
func NewConnection(server iface.IServer, conn *net.TCPConn, connID uint32, handler iface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:   server,
		Conn:        conn,
		ConnID:      connID,
		MsgHandle:   handler,
		isClosed:    false,
		ExitChan:    make(chan bool, 1),
		msgChan:     make(chan []byte, 1),
		msgBuffChan: make(chan []byte, utils.GlobalObject.MaxMsgChanLen), //不要忘记初始化
		property:    make(map[string]interface{}),
	}

	// 添加到ConnManager中
	c.TcpServer.GetConnMgr().Add(c)

	return c
}

// 启动链接
func (c *Connection) Start() {
	fmt.Println("Connection Start... ConnID = ", c.ConnID)
	// 启动从当前链接读数据的业务
	go c.StartReader()

	//启动写数据的业务
	go c.StartWriter()

	//按照用户传递进来的创建连接时需要处理的业务，执行钩子方法
	c.TcpServer.CallOnConnStart(c)
}

// 停止链接
func (c *Connection) Stop() {
	fmt.Println("Connection Stop... ConnID = ", c.ConnID)

	// 已关闭
	if c.isClosed {
		return
	}
	c.isClosed = true

	//如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用
	c.TcpServer.CallOnConnStop(c)

	c.Conn.Close()     // 关闭socket连接
	c.ExitChan <- true //告知Writer关闭

	// 从ConnManager中删除
	c.TcpServer.GetConnMgr().Remove(c)

	// 回收资源
	close(c.ExitChan)
	close(c.msgChan)
}

/**
* 读
**/
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine] is running...")
	defer fmt.Println("[Reader] is exit, Connection ID = ", c.ConnID, ", address is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		//headBuf := make([]byte, GetHeadLen())
		//if _, err := io.ReadFull(c.GetTCPConnect(), headBuf); err != nil {
		//	fmt.Println("Connection Read head error: ", err)
		//	break
		//}
		//
		//// 拆包，得到 MsgID 和 MsgData
		//msg, err := UnPack(headBuf)
		//if err != nil {
		//	fmt.Println("Message UnPack error: ", err)
		//	break
		//}
		//
		//// 根据DataLen，接着读取消息中的Data
		//if msg.GetDataLen() > 0 {
		//	dataBuf := make([]byte, msg.GetDataLen())
		//	if _, err := io.ReadFull(c.GetTCPConnect(), dataBuf); err != nil {
		//		fmt.Println("Connection Read Message data error: ", err)
		//		break
		//	}
		//
		//	msg.SetData(dataBuf)
		//}
		msg, err := UnPack(c.GetTCPConnect())
		if err != nil {
			fmt.Println("Connection UnPack error: ", err)
			break
		}

		// 得到当前conn数据的request请求数据
		req := &Requset{
			conn: c,
			msg:  msg,
		}

		// 执行注册的路由方法
		if utils.GlobalObject.WorkerPoolSize > 0 { // 已经开启工作池
			c.MsgHandle.AddRequestToTaskQueue(req) // 交给worker池
		} else { // 未开启工作池，直接处理
			go c.MsgHandle.DoMsgHandler(req)
		}
	}

}

/**
* 写消息Goroutine，专门发送给客户端消息
**/
func (c *Connection) StartWriter() {
	fmt.Println("[Connection]---(Writer Goroutine) is running...")
	defer fmt.Println("[Connection]---(Writer) is exit, Connection ID = ", c.ConnID, ", address is ", c.RemoteAddr().String())

	// 循环等待msgChan的消息，收到就发送给客户端
	for {
		select {
		case data := <-c.msgChan: //有数据，写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("[Connection]---send data error: ", err)
			}
		case data, ok := <-c.msgBuffChan: //针对有缓冲channel需要些的数据处理
			if ok {
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("[Connection]---send buff data error: ", err)
					return
				} else {
					fmt.Println("[Connection]---msgBuffChan is Closed")
					break
				}
			}
		case <-c.ExitChan: //Reader退出了，Writer也要退出
			return
		}
	}
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
	c.msgChan <- binaryMsg

	return nil
}

/**
* 直接将Message数据发送给远程的TCP客户端(有缓冲)
**/
func (c *Connection) SendBuffMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("[Connection] closed when send buff msg")
	}
	//将data封包，并且发送
	msg, err := Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Println("[Connection] Pack error msg id = ", msgId)
		return errors.New("[Connection] Pack error msg ")
	}

	//写回客户端
	c.msgBuffChan <- msg

	return nil
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

// 设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

// 获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

// 移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
