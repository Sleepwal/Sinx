# Sinx
学习[zinx](https://github.com/aceld/zinx)1.0开源框架

# 使用
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

## V0.3 基础router模块
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

## V0.4 全局配置模块
跟3.0相比，main函数中使用代码一样。

但是增加了全局配置，全局配置在conf/sinx.yaml中，可以配置一些全局参数，比如：
```yaml
host: 127.0.0.1
port: 8888
name: sinx demoServerApp
version: V0.4
maxConn: 3
```
