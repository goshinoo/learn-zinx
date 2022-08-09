package znet

import (
	"errors"
	"fmt"
	"github.com/goshinoo/learn-zinx/ziface"
	"sync"
)

// ConnManager 连接管理模块
type ConnManager struct {
	//管理的连接集合
	connections map[uint32]ziface.IConnection
	//保护连接集合的读写锁
	connLock sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (c *ConnManager) Add(conn ziface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	c.connections[conn.GetConnID()] = conn
	fmt.Println("connectionID = ", conn.GetConnID(), " add to ConnManager success:conn num = ", c.Len())
}

func (c *ConnManager) Remove(conn ziface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	delete(c.connections, conn.GetConnID())
	fmt.Println("connectionID = ", conn.GetConnID(), " remove success:conn num = ", c.Len())
}

func (c *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	c.connLock.RLock()
	defer c.connLock.RUnlock()

	if conn, ok := c.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

func (c *ConnManager) Len() int {
	return len(c.connections)
}

func (c *ConnManager) ClearConn() {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	for connID, conn := range c.connections {
		conn.Stop()
		delete(c.connections, connID)
	}

	fmt.Println("clear all connections successfully! conn num = ", c.Len())
}
