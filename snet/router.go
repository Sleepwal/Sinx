package snet

import "github.com/SleepWalker/sinx/iface"

/**
* 实现Router时，嵌入BaseRouter基类，根据需求重写方法
* BaseRouter的方法都为空，使用BaseRouter就不需要实现PreHandle、PostHandle
* 如果直接继承IRouter接口，需要实现所有方法。
**/
type BaseRouter struct{}

// 处理conn业务前的钩子方法Hook
func (br *BaseRouter) PreHandle(request iface.IRequest) {}

// 处理业务的主方法Hook
func (br *BaseRouter) Handle(request iface.IRequest) {}

// 处理业务后Hook
func (br *BaseRouter) PostHandle(request iface.IRequest) {}
