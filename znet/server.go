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

		//3.阻塞等待客户端连接,处理客户端连接业务
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept error:", err)
				continue
			}

			//已经与客户端建立连接,简单回显功能
			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("recv buf err", err)
						return
					}

					fmt.Printf("recv client :%s,cnt = %d\n", buf, cnt)

					//回显功能
					if _, err := conn.Write(buf[:cnt]); err != nil {
						fmt.Println("writer back buf err:", err)
						return
					}
				}
			}()
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

// NewServer 初始化Server模块的方法
func NewServer(name string) ziface.IServer {
	return &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
}
