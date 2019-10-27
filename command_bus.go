package cromberbus

import (
	"reflect"

	"github.com/chiguirez/cromberbus/v2/guard"
)

type Command interface{}

type CommandBus interface {
	Dispatch(command Command) error
}

type CommandHandler reflect.Value

func NewCommandHandler(handleFunc interface{}) CommandHandler {
	guard.CommandHandler(handleFunc)
	return CommandHandler(reflect.ValueOf(handleFunc))
}

func (c CommandHandler) Call(command Command) error {
	values := []reflect.Value{reflect.ValueOf(command)}
	res := reflect.Value(c).Call(values)
	err, ok := res[0].Interface().(error)
	if ok {
		return err
	}
	return nil
}

func (c CommandHandler) CommandName() string {
	return reflect.Value(c).Type().In(0).Name()
}
