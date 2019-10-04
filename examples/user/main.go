package main

import (
	"fmt"

	"github.com/chiguirez/cromberbus"
)

type RegisterUserCommand struct {
	email string
}

type RegisterUserCommandHandler struct{}

func (h RegisterUserCommandHandler) Handle(registerUserCommand RegisterUserCommand) error {
	fmt.Println("registering", registerUserCommand.email)
	return nil
}

type LoggingMiddleware struct{}

func (m LoggingMiddleware) Execute(command cromberbus.Command, next cromberbus.CommandCallable) error {
	fmt.Println("Execution of logging middleware")

	return next(command)
}

func main() {
	mapHandlerResolver := cromberbus.NewMapHandlerResolver()
	mapHandlerResolver.AddHandler(new(RegisterUserCommandHandler).Handle)
	bus := cromberbus.NewCromberBus(&mapHandlerResolver, new(LoggingMiddleware))
	command := RegisterUserCommand{"some@email.com"}
	err := bus.Dispatch(command)
	if err != nil {
		fmt.Printf("WARNING! an error occurred '%s'", err.Error())
	}
}
