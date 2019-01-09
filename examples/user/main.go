package main

import (
	"fmt"
	"github.com/hosseio/cromberbus"
	"github.com/pkg/errors"
)

type RegisterUserCommand struct {
	email string
}

type RegisterUserCommandHandler struct {}

func (h *RegisterUserCommandHandler) Handle(command cromberbus.Command) error {
	registerUserCommand, ok := command.(RegisterUserCommand)
	if !ok {
		return fmt.Errorf("Could not handle a non register user command")
	}

	fmt.Printf("%s was registered", registerUserCommand.email)
	return nil
}

func main() {
	mapHandlerResolver := cromberbus.NewMapHandlerResolver()
	mapHandlerResolver.AddHandler(new(RegisterUserCommand), new(RegisterUserCommandHandler))
	bus := cromberbus.NewCromberBus(&mapHandlerResolver)
	command := RegisterUserCommand{"some@email.com"}
	err := bus.Dispatch(command)
	if err != nil {
		fmt.Println(errors.Cause(err))
	}
}
