package maphandler

import (
	"fmt"

	"github.com/chiguirez/cromberbus/v3/commandhandler"
	"github.com/chiguirez/cromberbus/v3/maphandler/typer"
)

type Command interface{}

//go:generate moq -out command_handler_resolver_mock.go . CommandHandlerResolver
type CommandHandlerResolver interface {
	Resolve(command Command) (commandhandler.Class, error)
}

type Resolver struct {
	handlers map[string]commandhandler.Class
}

func NewResolver() Resolver {
	return Resolver{
		map[string]commandhandler.Class{},
	}
}

var ErrHandlerNotFound = fmt.Errorf("could not find command handler")

func (r Resolver) Resolve(command Command) (commandhandler.Class, error) {
	handler, ok := r.handlers[typer.Identify(command)]
	if !ok {
		return commandhandler.Class{}, ErrHandlerNotFound
	}

	return handler, nil
}

func (r Resolver) AddHandler(handlerFn commandhandler.Class) {
	r.handlers[handlerFn.CommandName()] = handlerFn
}
