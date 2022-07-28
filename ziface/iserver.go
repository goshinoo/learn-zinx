package ziface

// IServer 服务器接口
type IServer interface {
	// Start 启动服务器
	Start()
	// Stop 停止服务器
	Stop()
	// Serve 运行服务器
	Serve()
	// AddRouter 路由功能:给当前服务注册一个路由方法,供客户端的连接处理使用
	AddRouter(msgID uint32, router IRouter)
	// GetConnMgr 获取连接管理器
	GetConnMgr() IConnManager

	SetOnConnStart(func(connection IConnection))

	SetOnConnStop(func(connection IConnection))

	CallOnConnStart(connection IConnection)

	CallOnConnStop(connection IConnection)
}
