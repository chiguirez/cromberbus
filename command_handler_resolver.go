package cromberbus

import (
	"fmt"
	"reflect"
)

//go:generate moq -out command_handler_resolver_mock.go . CommandHandlerResolver
type CommandHandlerResolver interface {
	Resolve(command Command) (CommandHandler, error)
}

type MapHandlerResolver struct {
	handlers map[string]CommandHandler
}

func NewMapHandlerResolver() MapHandlerResolver {
	return MapHandlerResolver{
		map[string]CommandHandler{},
	}
}

func (r MapHandlerResolver) Resolve(command Command) (CommandHandler, error) {
	handler, ok := r.handlers[r.typeOf(command)]
	if !ok {
		return nil, fmt.Errorf("could not find command handler")
	}

	return handler, nil
}

func (r MapHandlerResolver) typeOf(command Command) string {
	if reflect.TypeOf(command).Kind() == reflect.Ptr {
		return reflect.TypeOf(command).Elem().String()
	}
	return reflect.TypeOf(command).String()
}

func (r MapHandlerResolver) AddHandler(command Command, handler CommandHandler) {
	r.handlers[r.typeOf(command)] = handler
}

