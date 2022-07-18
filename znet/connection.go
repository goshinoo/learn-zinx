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
	MsgHandler ziface.IMsgHandler
}

func NewConnection(conn *net.TCPConn, connID uint32, handler ziface.IMsgHandler) *Connection {
	return &Connection{
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		MsgHandler: handler,
		ExitChan:   make(chan bool, 1),
	}
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connId = ", c.ConnID, " Reader id exit,remote addr is ", c.RemoteAddr().String())
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

func (c *Connection) Start() {
	fmt.Println("Conn Start().. ConnID = ", c.ConnID)
	go c.StartReader()
}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop().. ConnID = ", c.ConnID)

	if c.isClosed {
		return
	}
	c.isClosed = true
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

	_, err = c.Conn.Write(binary)
	if err != nil {
		fmt.Println("write msg error", err)
		return err
	}
	return nil
}
