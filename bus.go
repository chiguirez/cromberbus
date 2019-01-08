package cromberbus

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
)

type Command interface {}

type CommandHandler interface {
	Handle(command Command) error
}

type CommandBus interface {
	Dispatch(command Command) error
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

func (r *MapHandlerResolver) typeOf(command Command) reflect.Type{
	if reflect.TypeOf(command).Kind() == reflect.Ptr {
		return reflect.TypeOf(command).Elem()
	}
	return reflect.TypeOf(command)
}

func (r *MapHandlerResolver) AddHandler(command Command, handler CommandHandler) {
	r.handlers[r.typeOf(command)] = handler
}

type CromberBus struct {
	handlerResolver HandlerResolver
}

func NewCromberBus(handlerResolver HandlerResolver) CromberBus {
	return CromberBus{handlerResolver}
}

func (b *CromberBus) Dispatch(command Command) error {
	handler, err := b.handlerResolver.Resolve(command)
	if err != nil {
		return errors.Wrap(err, "could not dispatch command")
	}

	return handler.Handle(command)
}
