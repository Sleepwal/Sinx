package iface

/**
* 消息管理抽象层
**/
type IMsgHandle interface {
	DoMsgHandler(request IRequest) // 执行对应的Router

	StartWorkerPool()            //启动worker工作池
	StartOneWorker(workerID int) //启动一个worker工作流程

	AddRouter(msgId uint32, router IRouter) // 为消息添加具体的处理逻辑
	AddRequestToTaskQueue(request IRequest) //将请求给TaskQueue，由worker进行处理
}
