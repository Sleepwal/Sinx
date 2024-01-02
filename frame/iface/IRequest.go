package iface

/*
* Request包中，封装：
* 客户端请求的链接 和 请求的数据
 */
type IRequest interface {
	// 链接
	GetConnection() IConnection

	// 请求的数据、ID
	GetData() []byte
	GetMsgID() uint32
}
