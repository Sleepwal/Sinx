package iface

type IServer interface {
	// 启动服务器
	Start()
	// 停止
	Stop()
	// 运行
	Serve()
}