package znet

import (
	"fmt"
	"net"
	"zinx/ziface"
)

// Server IServer的接口实现,定义Server的服务器模块
type Server struct {
	//服务器名称
	Name string
	//服务器绑定IP版本
	IPVersion string
	//服务器监听IP
	IP string
	//服务器监听端口
	Port int
	//当前的server添加router,server注册的连接对应的处理业务
	Router ziface.IRouter
}

func (s *Server) Start() {
	fmt.Printf("[Start] Server Listener at IP:%s,Port %d is starting\n", s.IP, s.Port)
	go func() {
		//1.获取TCP的addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error:", err)
			return
		}
		//2.监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err:", err)
			return
		}

		fmt.Println("start zinx server success,", s.Name, "Listening...")
		var cid uint32
		cid = 0

		//3.阻塞等待客户端连接,处理客户端连接业务
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept error:", err)
				continue
			}

			dealConn := NewConnection(conn, cid, s.Router)
			cid++

			//启动当前的连接业务
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	//TODO
}

func (s *Server) Serve() {
	//启动服务功能
	s.Start()

	//TODO 做启动服务后的额外功能

	//阻塞
	select {}
}

func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	fmt.Println("Add Router Succ!")
}

// NewServer 初始化Server模块的方法
func NewServer(name string) ziface.IServer {
	return &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
		Router:    nil,
	}
}
