package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/utils"
	"zinx/ziface"
)

type Connection struct {
	//当前Conn属于哪个server
	TcpServer  ziface.IServer
	Conn       *net.TCPConn
	ConnID     uint32
	isClosed   bool
	ExitChan   chan bool
	msgChan    chan []byte        //无缓冲通道,用于读写goroutine之间的消息通信
	MsgHandler ziface.IMsgHandler //消息的管理msgId和对应处理业务API的关系
	//连接属性
	property map[string]interface{}
	//保护连接属性的锁
	propertyLock sync.RWMutex
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if v, ok := c.property[key]; ok {
		return v, nil
	} else {
		return nil, errors.New("no property")
	}
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, handler ziface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		MsgHandler: handler,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
		property:   make(map[string]interface{}),
	}

	c.TcpServer.GetConnMgr().Add(c)

	return c
}

func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running...]")
	defer fmt.Println("[connId = ", c.ConnID, " ,Reader exit,remote addr is ", c.RemoteAddr().String(), "]")
	defer c.Stop()

	for {
		//创建拆包对象
		dp := NewDataPack()
		headData := make([]byte, dp.GetHeadLen())
		_, err := io.ReadFull(c.GetTCPConnection(), headData)
		if err != nil {
			fmt.Println("read msg head error", err)
			return
		}
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			return
		}

		data := make([]byte, msg.GetDataLen())
		if msg.GetDataLen() > 0 {
			_, err := io.ReadFull(c.GetTCPConnection(), data)
			if err != nil {
				fmt.Println("read msg data error", err)
				return
			}
		}

		msg.SetData(data)

		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			//执行注册的路由方法
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

//StartWriter 写消息的goroutine,专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running...]")
	defer fmt.Println("[connId = ", c.ConnID, " ,Writer exit,remote addr is ", c.RemoteAddr().String(), "]")
	defer close(c.msgChan)

	//不断阻塞等待channel消息
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error,", err)
				return
			}
		case <-c.ExitChan:
			//代表Reader已经退出,Writer也退出
			return
		}
	}
}

func (c *Connection) Start() {
	fmt.Println("Conn Start().. ConnID = ", c.ConnID)
	//启动从当前连接的读数据业务
	go c.StartReader()
	//启动从当前连接的写数据业务
	go c.StartWriter()

	c.TcpServer.CallOnConnStart(c)
}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop().. ConnID = ", c.ConnID)

	if c.isClosed {
		return
	}

	c.isClosed = true

	c.TcpServer.CallOnConnStop(c)

	c.TcpServer.GetConnMgr().Remove(c)
	//关闭socket连接
	_ = c.Conn.Close()
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("connection closed")
	}

	dp := NewDataPack()
	binary, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error ", err)
		return err
	}

	c.msgChan <- binary
	return nil
}
