package iface

/**
* 消息管理抽象层
**/
type IMsgHandle interface {
	// 执行对应的Router
	DoMsgHandler(request IRequest)

	// 为消息添加具体的处理逻辑
	AddRouter(msgId uint32, router IRouter)
}
