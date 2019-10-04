package cromberbus

import (
	"fmt"

	"github.com/chiguirez/cromberbus/typer"
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
	handler, ok := r.handlers[typer.Identify(command)]
	if !ok {
		return CommandHandler{}, fmt.Errorf("could not find command handler")
	}

	return handler, nil
}

func (r MapHandlerResolver) AddHandler(handlerFn interface{}) {
	commandHandler := NewCommandHandler(handlerFn)
	r.handlers[commandHandler.CommandName()] = commandHandler
}
