package net

import (
	"fmt"
	"sync"

	"SleepXLink/iface"
)

/**
* 连接管理模块
**/
type ConnManager struct {
	connections map[uint32]iface.IConnection //管理的连接集合
	connLock    sync.RWMutex                 //保护连接集合的读、写锁
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]iface.IConnection),
	}
}

// 添加连接
func (cm *ConnManager) Add(conn iface.IConnection) {
	// 保护共享资源map，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// conn加入connections中
	cm.connections[conn.GetConnID()] = conn
	fmt.Println("[ConnManager] Add conn success, connID:", conn.GetConnID(),
		", conn num = ", cm.Len())
}

// 删除连接
func (cm *ConnManager) Remove(conn iface.IConnection) {
	// 保护共享资源map，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 删除conn
	delete(cm.connections, conn.GetConnID())
	fmt.Println("[ConnManager] Remove conn success, connID:", conn.GetConnID(),
		", conn num = ", cm.Len())
}

// 根据conn获取连接
func (cm *ConnManager) Get(connID uint32) (iface.IConnection, error) {
	// 保护共享资源map，加读锁
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	if conn, ok := cm.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, fmt.Errorf("[ConnManager] connection not found, connID = %d", connID)
	}
}

// 获取连接总数
func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

// 清除并停止所有连接
func (cm *ConnManager) ClearConn() {
	// 保护共享资源map，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 停止并删除全部的conn
	for connID, conn := range cm.connections {
		conn.Stop()                    // 停止
		delete(cm.connections, connID) // 删除
	}
	fmt.Println("[ConnManager] Clear All connections success! conn num = ", cm.Len())
}
