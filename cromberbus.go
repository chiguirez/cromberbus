package cromberbus

type CromberBus struct {
	middlewares MiddlewareList
}

func NewCromberBus(handlerResolver CommandHandlerResolver, middlewares ...Middleware) CromberBus {
	commandHandlingMiddleware := commandHandlingMiddleware{handlerResolver}
	middlewareList := NewMiddlewareList(commandHandlingMiddleware).Queue(middlewares...)

	return CromberBus{middlewareList}
}

func (b CromberBus) Dispatch(command Command) {
	b.middlewares.start(command)
}
