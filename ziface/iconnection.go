package ziface

import "net"

type IConnection interface {
	// Start 启动连接 让当前连接准备开始工作
	Start()
	// Stop 停止链接 结束当前链接的工作
	Stop()
	// GetTCPConnection 获取当前链接绑定的socket conn
	GetTCPConnection() *net.TCPConn
	// GetConnID 获取当前链接模块的ID
	GetConnID() uint32
	// RemoteAddr 获取远程客户端的TCP状态 IP PORT
	RemoteAddr() net.Addr
	// SendMsg 发送数据 将数据发送给远程的客户端
	SendMsg(uint32, []byte) error
	// SetProperty 设置连接属性
	SetProperty(key string, value interface{})
	// GetProperty 获取连接属性
	GetProperty(key string) (interface{}, error)
	// RemoveProperty 移除连接属性
	RemoveProperty(key string)
}

// HandleFunc 处理连接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error
