package znet

import (
	"fmt"
	"net"
	"zinx/utils"
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
	//当前的server的消息管理模块,绑定MsgId和对应处理业务API关系\
	MsgHandler ziface.IMsgHandler
	//该server的连接管理器
	ConnMgr ziface.IConnManager
	//该server创建连接之后自动调用的Hook函数
	OnConnStart func(conn ziface.IConnection)
	//该server关闭连接之前自动调用的Hook函数
	OnConnStop func(conn ziface.IConnection)
}

func (s *Server) SetOnConnStart(f func(connection ziface.IConnection)) {
	s.OnConnStart = f
}

func (s *Server) SetOnConnStop(f func(connection ziface.IConnection)) {
	s.OnConnStop = f
}

func (s *Server) CallOnConnStart(connection ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("Call OnConnStart()...")
		s.OnConnStart(connection)
	}
}

func (s *Server) CallOnConnStop(connection ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("Call OnConnStop()...")
		s.OnConnStop(connection)
	}
}

func (s *Server) Start() {
	fmt.Printf("[Zinx]Server Name:%s,listener at Ip:%s,Port:%d is starting...\n", utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx]Version:%s,MaxConn:%d,MaxPackageSize:%d\n", utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)

	go func() {
		//0.开启消息队列
		s.MsgHandler.StartWorkerPool()

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

			//设置最大连接个数的判断,如果超过最大连接数,则关闭此链接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				//todo 给客户端响应错误
				fmt.Println("too many conns,maxConn = ", utils.GlobalObject.MaxConn)
				_ = conn.Close()
				continue
			}

			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			//启动当前的连接业务
			go dealConn.Start()
		}
	}()
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server name ", s.Name)
	s.ConnMgr.ClearConn()
}

func (s *Server) Serve() {
	//启动服务功能
	s.Start()

	//TODO 做启动服务后的额外功能

	//阻塞
	select {}
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router Succ!")
}

// NewServer 初始化Server模块的方法
func NewServer() ziface.IServer {
	return &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandler(),
		ConnMgr:    NewConnManager(),
	}
}
