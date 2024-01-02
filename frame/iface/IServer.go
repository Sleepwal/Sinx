package iface

type IServer interface {
	// 启动服务器
	Start()
	// 停止
	Stop()
	// 运行
	Serve()
	//注册一个路由
	AddRouter(msgID uint32, router IRouter)
	// 获取当前Server的连接管理器
	GetConnMgr() IConnManager

	//设置该Server的连接创建时Hook函数
	SetOnConnStart(func(IConnection))
	//设置该Server的连接断开时的Hook函数
	SetOnConnStop(func(IConnection))
	//调用连接OnConnStart Hook函数
	CallOnConnStart(conn IConnection)
	//调用连接OnConnStop Hook函数
	CallOnConnStop(conn IConnection)
}
