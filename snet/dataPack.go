package snet

import (
	"bytes"
	"encoding/binary"
	"errors"

	"sinx/iface"
	"sinx/utils"
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
func UnPack(binaryData []byte) (iface.IMessage, error) {
	reader := bytes.NewReader(binaryData)
	msg := &Message{}

	// 1. 读取长度
	err := binary.Read(reader, binary.LittleEndian, &msg.DataLen)
	if err != nil {
		return nil, err
	}

	// 2. 读取ID
	err = binary.Read(reader, binary.LittleEndian, &msg.ID)
	if err != nil {
		return nil, err
	}

	// 判断数据包 是否超过 最大长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.GetDataLen() > utils.GlobalObject.MaxPackageSize {
		return msg, errors.New("package size over limit")
	}

	//这里只需要把head的数据拆包出来就可以了，然后再通过head的长度，再从conn读取一次数据
	return msg, nil
}
