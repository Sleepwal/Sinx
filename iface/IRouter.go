package iface

/**
* 路由接口
* 数据都是IRequest
**/
type IRouter interface {
	// 处理conn业务前的钩子方法Hook
	PreHandle(request IRequest)

	// 处理业务的主方法Hook
	Handle(request IRequest)

	// 处理业务后Hook
	PostHandle(request IRequest)
}
