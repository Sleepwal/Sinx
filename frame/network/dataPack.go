package network

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"SleepXLink/iface"
	"SleepXLink/utils"
)

/**
* 封包、拆包 工具
* 面向TCP连接中的数据流，处理TCP粘包问题
**/

// 获取包的头的长度
func GetHeadLen() uint32 {
	// DataLen uint32(4字节) + ID uint32(4字节)
	return 8
}

// 封包 len|id|data
func Pack(msg iface.IMessage) ([]byte, error) {
	// 存放字节的缓冲
	dataBuf := bytes.NewBuffer([]byte{})

	// 按照包格式写入
	// 1. Len数据长度
	err := binary.Write(dataBuf, binary.LittleEndian, msg.GetDataLen())
	if err != nil {
		return nil, err
	}

	// 2. MsgID
	err = binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgId())
	if err != nil {
		return nil, err
	}

	// 3. data数据
	err = binary.Write(dataBuf, binary.LittleEndian, msg.GetData())
	if err != nil {
		return nil, err
	}

	return dataBuf.Bytes(), nil
}

// 拆包 len|id|data
func UnPack(reader io.Reader) (iface.IMessage, error) {
	//================1.head头部部分================
	//1.1 接收客户端消息的head，8个字节
	headData := make([]byte, GetHeadLen())
	if _, err := io.ReadFull(reader, headData); err != nil {
		fmt.Println("Connection Read head error: ", err)
		return nil, err
	}

	//1.2 读取head信息
	bytesReader := bytes.NewReader(headData)
	msg := &Message{}

	err := binary.Read(bytesReader, binary.LittleEndian, &msg.DataLen) //消息长度
	if err != nil {
		return nil, err
	}

	err = binary.Read(bytesReader, binary.LittleEndian, &msg.ID) //消息ID
	if err != nil {
		return nil, err
	}

	//================2.data数据部分================
	//2.1 判断数据包 是否超过 最大长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.GetDataLen() > utils.GlobalObject.MaxPackageSize {
		return msg, errors.New("package size over limit")
	}

	//2.2 根据DataLen，接着读取消息中的Data
	if msg.GetDataLen() > 0 {
		dataBuf := make([]byte, msg.GetDataLen())
		if _, err := io.ReadFull(reader, dataBuf); err != nil {
			fmt.Println("Connection Read Message data error: ", err)
			return nil, err
		}

		msg.SetData(dataBuf)
	}

	return msg, nil
}
