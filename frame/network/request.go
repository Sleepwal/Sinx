package network

import "SleepXLink/iface"

type Requset struct {
	// 与客户端 建立好的链接
	conn iface.IConnection

	// 请求的数据
	msg iface.IMessage
}

// 链接
func (r *Requset) GetConnection() iface.IConnection {
	return r.conn
}

// 获取请求的数据
func (r *Requset) GetData() []byte {
	return r.msg.GetData()
}

// 请求的ID
func (r *Requset) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}
