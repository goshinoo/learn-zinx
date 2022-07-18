package ziface

// IMessage 将请求的消息封装到message中
type IMessage interface {
	GetMsgId() uint32
	GetDataLen() uint32
	GetData() []byte

	SetMsgId(uint32)
	SetData([]byte)
	SetDataLen(uint32)
}
