package cromberbus

import (
	"fmt"
	"reflect"
)

type Command interface{}

type CommandHandler interface {
	Handle(command Command) error
}

//TODO: what about error propagation and error handling
type CommandBus interface {
	Dispatch(command Command)
}

//go:generate moq -out handler_resolver_mock.go . HandlerResolver
type HandlerResolver interface {
	Resolve(command Command) (CommandHandler, error)
}

type MapHandlerResolver struct {
	handlers map[reflect.Type]CommandHandler
}

func NewMapHandlerResolver() MapHandlerResolver {
	return MapHandlerResolver{
		map[reflect.Type]CommandHandler{},
	}
}

func (r *MapHandlerResolver) Resolve(command Command) (CommandHandler, error) {
	handler, ok := r.handlers[r.typeOf(command)]
	if !ok {
		return nil, fmt.Errorf("could not find command handler")
	}

	return handler, nil
}

func (r *MapHandlerResolver) typeOf(command Command) reflect.Type {
	if reflect.TypeOf(command).Kind() == reflect.Ptr {
		return reflect.TypeOf(command).Elem()
	}
	return reflect.TypeOf(command)
}

func (r *MapHandlerResolver) AddHandler(command Command, handler CommandHandler) {
	r.handlers[r.typeOf(command)] = handler
}

type CommandCallable func(command Command)

type Middleware interface {
	Execute(command Command, next CommandCallable)
}

type commandHandlingMiddleware struct {
	handlerResolver HandlerResolver
}

func (m *commandHandlingMiddleware) Execute(command Command, next CommandCallable) {
	//TODO: why are we ignoring this error?
	handler, _ := m.handlerResolver.Resolve(command)

	if err := handler.Handle(command); err!=nil{
		return
	}

	next(command)
}

type MiddlewareList []Middleware

func NewMiddlewareList(commandHandler *commandHandlingMiddleware) MiddlewareList {
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

type CromberBus struct {
	middlewares MiddlewareList
}

func NewCromberBus(handlerResolver HandlerResolver, middlewares ...Middleware) CromberBus {
	commandHandlingMiddleware := &commandHandlingMiddleware{handlerResolver}
	middlewareList := NewMiddlewareList(commandHandlingMiddleware).Queue(middlewares...)

	return CromberBus{middlewareList}
}

func (b *CromberBus) Dispatch(command Command) {
	b.middlewares.start(command)
}
