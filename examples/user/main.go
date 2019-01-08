package main

import (
	"bitbucket.org/hosseio/cromberbus"
	"fmt"
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
	registerUserCommand := RegisterUserCommand{"some@email.com"}
	registerUserCommandHandler := &RegisterUserCommandHandler{}
	mapHandlerResolver.AddHandler(registerUserCommand, registerUserCommandHandler)

	bus := cromberbus.NewCromberBus(&mapHandlerResolver)
	err := bus.Dispatch(registerUserCommand)
	if err != nil {
		fmt.Println(errors.Cause(err))
	}
}
