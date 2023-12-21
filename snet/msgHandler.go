package snet

import (
	"fmt"
	"strconv"

	"sinx/iface"
	"sinx/utils"
)

/**
* 消息模块模块
**/
type MsgHandle struct {
	Apis           map[uint32]iface.IRouter // 存放每个MsgID对应的处理逻辑
	TaskQueue      []chan iface.IRequest    // 负责worker取任务的消息队列
	WorkerPoolSize uint32                   // 业务工作worker池的worker数量
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]iface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,                              //全局配置中获取
		TaskQueue:      make([]chan iface.IRequest, utils.GlobalObject.WorkerPoolSize), //一个worker对应一个queus
	}
}

// 执行对应的Router
func (mh *MsgHandle) DoMsgHandler(request iface.IRequest) {
	// 1.从Request获取MsgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgId = ", request.GetMsgID(), " is not FOUND!")
	}

	// 2.调用MsgID对应的Router业务
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgId uint32, router iface.IRouter) {
	// 1.判断当前msg是否已绑定
	if _, ok := mh.Apis[msgId]; ok {
		panic("msgId is exist, id = " + strconv.Itoa(int(msgId)))
	}

	// 2.msg与API绑定
	mh.Apis[msgId] = router
	fmt.Println("add api success, msgId = ", msgId)
}

/**
* 启动worker工作池
* 开启工作池的动作只能发生一次
**/
func (mh *MsgHandle) StartWorkerPool() {
	// 根据workerPoolSize分别开启worker，每个worker开启一个goroutine
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 1.当前worker对应的channel消息队列，开辟空间
		mh.TaskQueue[i] = make(chan iface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 2.启动当前worker，阻塞等待消息 从channel传递进来
		go mh.StartOneWorker(i)
	}
}

/**
* 启动一个worker工作流程
**/
func (mh *MsgHandle) StartOneWorker(workerID int) {
	fmt.Println("[MsgHandler] Worker ID = ", workerID, " is started..")

	// 阻塞等待channel的消息
	for {
		select {
		case request := <-mh.TaskQueue[workerID]: //有消息
			mh.DoMsgHandler(request) // 执行当前Request的MsgID绑定的业务
		}
	}
}

/**
* 将请求给TaskQueue，由worker进行处理
**/
func (mh *MsgHandle) AddRequestToTaskQueue(request iface.IRequest) {
	// 1.将请求平均分配给worker，根据ConnID进行分配
	workerID := request.GetConnection().GetConnID() % utils.GlobalObject.WorkerPoolSize
	fmt.Println("Add ConnId = ", request.GetConnection().GetConnID(),
		", request MsgID = ", request.GetMsgID(),
		" to WorkerID = ", workerID)

	// 2.发送给worker对应的TaskQueue
	mh.TaskQueue[workerID] <- request
}
