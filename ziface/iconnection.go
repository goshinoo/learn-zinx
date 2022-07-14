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
	// Send 发送数据 将数据发送给远程的客户端
	Send(data []byte) error
}

// HandleFunc 处理连接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error
