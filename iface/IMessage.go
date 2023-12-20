package iface

/**
* 请求的消息封装到Message中，定义抽象接口
**/

type IMessage interface {
	//===Getter===
	GetMsgId() uint32
	GetDataLen() uint32
	GetData() []byte

	//===Setter===
	SetMsgId(msgId uint32)
	SetDataLen(msgLen uint32)
	SetData(msgData []byte)
}
