package ziface

type IRequest interface {
	// GetConnection 得到当前连接
	GetConnection() IConnection
	// GetData 得到请求的数据
	GetData() []byte
	GetMsgId() uint32
}
