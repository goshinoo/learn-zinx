package znet

import (
	"zinx/ziface"
)

// BaseRouter 实现router时,先嵌入BaseRouter基类,根据需要对这个基类的方法进行重写
type BaseRouter struct {
}

/*这里之所以方法都为空,因为有的Router不希望有钩子业务*/

func (b *BaseRouter) PreHandle(request ziface.IRequest) {}

func (b *BaseRouter) Handle(request ziface.IRequest) {}

func (b *BaseRouter) PostHandle(request ziface.IRequest) {}
