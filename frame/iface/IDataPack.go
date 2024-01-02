package iface

import "io"

/****************************************
@Author : SleepWalker
@Description: 封包数据和拆包数据，为传输数据添加头部信息，用于处理TCP粘包问题。
@Time : 2024/1/2 16:30
****************************************/

type IDataPack interface {
	GetHeadLen() uint32                //获取包头长度方法
	Pack(msg IMessage) ([]byte, error) //封包方法
	Unpack([]byte) (io.Reader, error)  //拆包方法
}
