# Sinx
学习[zinx](https://github.com/aceld/zinx)1.0开源框架

# V1.0 链接管理和连接属性设置

# 使用该框架开发

配置文件
```yaml
host: 127.0.0.1
port: 8888
name: Sinx V1.0 demoServerApp
maxConn: 3
workerPoolSize: 10
```

server.go
```go
func main() {
	server := snet.NewServer()

	//注册链接hook回调函数
	server.SetOnConnStart(DoConnectionBegin)
	server.SetOnConnStop(DoConnectionLast)

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

// 创建连接的时候执行
func DoConnectionBegin(conn iface.IConnection) {
	fmt.Println("DoConnecionBegin is Called ... ")

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

	fmt.Println("DoConneciotnLost is Called ... ")
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