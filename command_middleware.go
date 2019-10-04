package cromberbus

type CommandCallable func(command Command) error

type Middleware interface {
	Execute(command Command, next CommandCallable) error
}

type commandHandlingMiddleware struct {
	handlerResolver CommandHandlerResolver
}

func (m commandHandlingMiddleware) Execute(command Command, next CommandCallable) error {
	handler, err := m.handlerResolver.Resolve(command)
	if err != nil {
		return err
	}

	if err := handler.Call(command); err != nil {
		return err
	}

	return next(command)
}

type MiddlewareList []Middleware

func NewMiddlewareList(commandHandler Middleware) MiddlewareList {
	return []Middleware{commandHandler}
}

func (m MiddlewareList) Queue(middleware ...Middleware) MiddlewareList {
	return append(m, middleware...)
}

func (m MiddlewareList) start(command Command) error {
	return m.getCallable(0)(command)
}

func (m MiddlewareList) lastIndex() int {
	return len(m) - 1
}

func (m MiddlewareList) getCallable(index int) CommandCallable {
	lastCallable := func(command Command) error { return nil }
	if index > m.lastIndex() {
		return lastCallable
	}

	return func(command Command) error {
		middleware := m[index]

		return middleware.Execute(command, m.getCallable(index+1))
	}
}
