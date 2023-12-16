package snet

import "github.com/SleepWalker/sinx/iface"

type Requset struct {
	// 与客户端 建立好的链接
	conn iface.IConnection

	// 请求的数据
	data []byte
}

// 链接
func (r *Requset) GetConnection() iface.IConnection {
	return r.conn
}

// 请求的数据
func (r *Requset) GetData() []byte {
	return r.data
}
