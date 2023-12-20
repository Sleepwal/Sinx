package snet

type Message struct {
	ID      uint32 //消息ID
	DataLen uint32 //消息长度
	Data    []byte //消息内容
}

// ===New===
func NewMessage(msgId uint32, data []byte) *Message {
	return &Message{
		ID:      msgId,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

// ===Getter===
func (m *Message) GetMsgId() uint32 {
	return m.ID
}
func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}
func (m *Message) GetData() []byte {
	return m.Data
}

// ===Setter===
func (m *Message) SetMsgId(msgId uint32) {
	m.ID = msgId
}
func (m *Message) SetDataLen(msgLen uint32) {
	m.DataLen = msgLen
}
func (m *Message) SetData(msgData []byte) {
	m.Data = msgData
}
