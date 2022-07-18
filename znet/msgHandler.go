package znet

import (
	"fmt"
	"zinx/ziface"
)

type MsgHandler struct {
	Api map[uint32]ziface.IRouter
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{Api: make(map[uint32]ziface.IRouter)}
}

func (m *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	router, ok := m.Api[request.GetMsgId()]
	if !ok {
		fmt.Println("route not found,msgId = ", request.GetMsgId())
		return
	}
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
}

func (m *MsgHandler) AddRouter(msgID uint32, router ziface.IRouter) {
	if _, ok := m.Api[msgID]; ok {
		fmt.Println("repeat api,msgId = ", msgID)
		return
	}
	m.Api[msgID] = router
	fmt.Println("Add api MsgID = ", msgID, " succ!")
}
