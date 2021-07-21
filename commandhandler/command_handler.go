package commandhandler

import (
	"context"
	"reflect"

	"github.com/chiguirez/cromberbus/v3/commandhandler/guard"
)

type Command interface{}

type Class struct {
	value reflect.Value
}

// New CommandHandler.
func New(handleFunc interface{}) Class {
	guard.CommandHandler(handleFunc)

	return Class{reflect.ValueOf(handleFunc)}
}

func (c Class) Call(ctx context.Context, command Command) error {
	values := []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(command)}
	res := c.value.Call(values)

	err, ok := res[0].Interface().(error)
	if ok {
		return err
	}

	return nil
}

func (c Class) CommandName() string {
	return c.value.Type().In(1).Name()
}
