# Sinx
学习[zinx](https://github.com/aceld/zinx)1.0开源框架

# 使用该框架开发
## V0.1-基础Server模块

```go
func main() {
	server := snet.NewServer("V0.1")
	server.Serve()
}
```

## V0.2-简单的链接封装和业务绑定

```go
func main() {
	server := snet.NewServer("V0.2")
	server.Serve()
}
```

## V0.3-基础router模块
```go
func main() {
	server := snet.NewServer("V0.3")
	//添加自定义router
	server.AddRouter(&PingRouter{})
	server.Serve()
}

// ping 自定义路由
type PingRouter struct {
	snet.BaseRouter
}

func (pr *PingRouter) PreHandle(request iface.IRequest) {
	fmt.Println("Call Router PreHandle...")
	_, err := request.GetConnection().GetTCPConnect().Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("cal back pre handle error: ", err)
	}
}

func (pr *PingRouter) Handle(request iface.IRequest) {
	fmt.Println("Call Router Handle...")
	_, err := request.GetConnection().GetTCPConnect().Write([]byte("ping ping...\n"))
	if err != nil {
		fmt.Println("cal back handle error: ", err)
	}
}

func (pr *PingRouter) PostHandle(request iface.IRequest) {
	fmt.Println("Call Router PostHandle...")
	_, err := request.GetConnection().GetTCPConnect().Write([]byte("after ping...\n"))
	if err != nil {
		fmt.Println("cal back post handle error: ", err)
	}
}
```

## V0.4-全局配置模块
- 跟3.0相比，main函数中使用代码一样。
- 但是增加了全局配置，全局配置在conf/sinx.yaml中，可以配置一些全局参数，比如：

```yaml
host: 127.0.0.1
port: 8888
name: sinx demoServerApp
version: V0.4
maxConn: 3
```

## V0.5 Message消息模块
- 解决TCP粘包问题
- 自定义消息包 【DataLen | MsgID | MsgData】
- 定义对应的封包、拆包方法

server.go
```go
func main() {
	server := snet.NewServer()
	//添加自定义router
	server.AddRouter(&PingRouter{})
	server.Serve()
}

// ping 自定义路由
type PingRouter struct {
	snet.BaseRouter
}

func (pr *PingRouter) Handle(request iface.IRequest) {
	fmt.Println("Call Router Handle...")
	// 1.读取客户端数据
	fmt.Println("recv from client: msgID = ", request.GetMsgID(), ", data = ", string(request.GetData()))

	// 2.回显ping
	err := request.GetConnection().SendMsg(request.GetMsgID(), []byte("ping...ping...\n"))
	if err != nil {
		fmt.Println("cal back handle error: ", err)
	}
}
```

client.go
```go
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
		msg, _ := snet.Pack(snet.NewMessage(uint32(i), []byte("Sinx V0.5 Client Test Message")))
		_, err := conn.Write(msg)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		//先读出流中的head部分
		headData := make([]byte, snet.GetHeadLen())
		_, err = io.ReadFull(conn, headData) //ReadFull 会把msg填充满为止
		if err != nil {
			fmt.Println("read head error")
			break
		}
		//将headData字节流 拆包到msg中
		recvMsg, err := snet.UnPack(headData)
		if err != nil {
			fmt.Println("server unpack err:", err)
			return
		}

		if recvMsg.GetDataLen() > 0 {
			//msg 是有data数据的，需要再次读取data数据
			recvMsg.SetData(make([]byte, recvMsg.GetDataLen()))

			//根据dataLen从io中读取字节流
			_, err := io.ReadFull(conn, recvMsg.GetData())
			if err != nil {
				fmt.Println("server unpack data err:", err)
				return
			}

			fmt.Println("==> Recv Msg: ID=", recvMsg.GetMsgId(),
				", len=", recvMsg.GetDataLen(),
				", data=", string(recvMsg.GetData()))
		}

		time.Sleep(1 * time.Second)
	}
}
```

## V0.6-多路由模式 

server.go
```go
func main() {
	server := snet.NewServer()
	//添加自定义router
	server.AddRouter(0, &PingRouter{})
	server.AddRouter(1, &HelloRouter{})
	server.Serve()
}

// ping 自定义路由
type PingRouter struct {
	snet.BaseRouter
}

func (pr *PingRouter) Handle(request iface.IRequest) {
	fmt.Println("Call Ping Router Handle...")
	// 1.读取客户端数据
	fmt.Println("recv from client: msgID = ", request.GetMsgID(), ", data = ", string(request.GetData()))

	// 2.回显ping
	err := request.GetConnection().SendMsg(200, []byte("ping...ping...\n"))
	if err != nil {
		fmt.Println("cal back handle error: ", err)
	}
}

// hello 自定义路由
type HelloRouter struct {
	snet.BaseRouter
}

func (hr *HelloRouter) Handle(request iface.IRequest) {
	fmt.Println("Call Hello Router Handle...")
	// 1.读取客户端数据
	fmt.Println("recv from client: msgID = ", request.GetMsgID(), ", data = ", string(request.GetData()))

	// 2.回显ping
	err := request.GetConnection().SendMsg(500, []byte("Hello Sinx!\n"))
	if err != nil {
		fmt.Println("cal back handle error: ", err)
	}
}
```



client.go
```go
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
		msg, _ := snet.Pack(snet.NewMessage(uint32(i)%2, []byte(fmt.Sprintf("Sinx Client%d Test Message", i))))
		_, err := conn.Write(msg)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		//先读出流中的head部分
		headData := make([]byte, snet.GetHeadLen())
		_, err = io.ReadFull(conn, headData) //ReadFull 会把msg填充满为止
		if err != nil {
			fmt.Println("read head error")
			break
		}
		//将headData字节流 拆包到msg中
		recvMsg, err := snet.UnPack(headData)
		if err != nil {
			fmt.Println("server unpack err:", err)
			return
		}

		if recvMsg.GetDataLen() > 0 {
			//msg 是有data数据的，需要再次读取data数据
			recvMsg.SetData(make([]byte, recvMsg.GetDataLen()))

			//根据dataLen从io中读取字节流
			_, err := io.ReadFull(conn, recvMsg.GetData())
			if err != nil {
				fmt.Println("server unpack data err:", err)
				return
			}

			fmt.Println("==> Recv Msg: ID=", recvMsg.GetMsgId(),
				", len=", recvMsg.GetDataLen(),
				", data=", string(recvMsg.GetData()))
		}

		time.Sleep(1 * time.Second)
	}
}

```

## V0.7
- Connection中读、写分离
  - 添加一个管道msgChan，读、写Goroutine进行通信

同V0.6

配置文件
```yaml
host: 127.0.0.1
port: 8888
name: Sinx V0.7 demoServerApp
maxConn: 3
```