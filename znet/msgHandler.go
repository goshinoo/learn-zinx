package znet

import (
	"fmt"
	"zinx/utils"
	"zinx/ziface"
)

type MsgHandler struct {
	//存放每个MsgID所对应的处理方法
	Api map[uint32]ziface.IRouter
	//负责worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
	//业务工作Worker池的worker数量
	WorkerPoolSize uint32
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Api:            make(map[uint32]ziface.IRouter),
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
	}
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

// StartWorkerPool 启动一个worker工作池
func (m *MsgHandler) StartWorkerPool() {
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		m.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskSize)
		go m.startOneWorker(i)
	}
}

//启动一个worker工作流程
func (m *MsgHandler) startOneWorker(workerID int) {
	fmt.Println("WorkerID = ", workerID, " is started...")
	for request := range m.TaskQueue[workerID] {
		m.DoMsgHandler(request)
	}
}

// SendMsgToTaskQueue 将消息交给Taskqueue,由worker进行处理
func (m *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	//将消息平均分配给不同的worker
	//根据客户端建立的connid分配
	workerId := request.GetConnection().GetConnID() % m.WorkerPoolSize
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(), "request MsgID = ", request.GetMsgId(), " to WorkerId = ", workerId)

	//将消息发送给对应worker的taskqueue
	m.TaskQueue[workerId] <- request
}
