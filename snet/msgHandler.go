package snet

import (
	"fmt"
	"strconv"

	"github.com/SleepWalker/sinx/iface"
)

/**
* 消息模块模块
**/
type MsgHandle struct {
	// 存放每个MsgID对应的处理逻辑
	Apis map[uint32]iface.IRouter
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]iface.IRouter),
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
