package cromberbus

type CommandCallable func(command Command)

type Middleware interface {
	Execute(command Command, next CommandCallable)
}

type commandHandlingMiddleware struct {
	handlerResolver CommandHandlerResolver
}

func (m commandHandlingMiddleware) Execute(command Command, next CommandCallable) {
	//TODO: why are we ignoring this error?
	handler, _ := m.handlerResolver.Resolve(command)

	if err := handler.Handle(command); err!=nil{
		return
	}

	next(command)
}

type MiddlewareList []Middleware

func NewMiddlewareList(commandHandler commandHandlingMiddleware) MiddlewareList {
	return []Middleware{commandHandler}
}

func (m MiddlewareList) Queue(middleware ...Middleware) MiddlewareList {
	return append(m, middleware...)
}

func (m MiddlewareList) start(command Command) {
	m.getCallable(0)(command)
}

func (m MiddlewareList) lastIndex() int {
	return len(m) - 1
}

func (m MiddlewareList) getCallable(index int) CommandCallable {
	lastCallable := func(command Command) {}
	if index > m.lastIndex() {
		return lastCallable
	}

	return func(command Command) {
		middleware := m[index]
		middleware.Execute(command, m.getCallable(index+1))
	}
}
