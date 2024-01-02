package iface

/**
* 连接管理模块抽象层
**/

type IConnManager interface {
	Add(conn IConnection)                   // 添加连接
	Remove(conn IConnection)                // 删除连接
	Get(connID uint32) (IConnection, error) // 根据conn获取连接
	Len() int                               // 获取连接总数
	ClearConn()                             // 清除并停止所有连接
}
