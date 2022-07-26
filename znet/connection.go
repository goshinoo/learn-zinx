package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/ziface"
)

type Connection struct {
	Conn       *net.TCPConn
	ConnID     uint32
	isClosed   bool
	ExitChan   chan bool
	msgChan    chan []byte        //无缓冲通道,用于读写goroutine之间的消息通信
	MsgHandler ziface.IMsgHandler //消息的管理msgId和对应处理业务API的关系
}

func NewConnection(conn *net.TCPConn, connID uint32, handler ziface.IMsgHandler) *Connection {
	return &Connection{
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		MsgHandler: handler,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
	}
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

		//执行注册的路由方法
		go c.MsgHandler.DoMsgHandler(&req)
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
}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop().. ConnID = ", c.ConnID)

	if c.isClosed {
		return
	}

	c.isClosed = true
	//关闭socket连接
	c.Conn.Close()

	close(c.ExitChan)
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
