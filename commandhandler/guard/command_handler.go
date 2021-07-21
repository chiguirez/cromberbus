package guard

import (
	"context"
	"reflect"
)

func FuncHaveExactlyTwoAttrib(handler interface{}) bool {
	return reflect.TypeOf(handler).NumIn() == 2
}

func FuncAttributeMustBeContext(handler interface{}) bool {
	ctxInterface := reflect.TypeOf((*context.Context)(nil)).Elem()

	return reflect.TypeOf(handler).In(0).Implements(ctxInterface)
}

func FuncAttributeMustBeStruct(handler interface{}) bool {
	return reflect.TypeOf(handler).In(1).Kind() == reflect.Struct
}

func FuncReturnsError(handler interface{}) bool {
	errInterface := reflect.TypeOf((*error)(nil)).Elem()

	if reflect.TypeOf(handler).NumOut() != 1 {
		return false
	}

	return reflect.TypeOf(handler).Out(0).Implements(errInterface)
}

func IsAFunc(handler interface{}) bool {
	return reflect.TypeOf(handler).Kind() == reflect.Func
}

func CommandHandler(handler interface{}) {
	if !IsAFunc(handler) {
		panic("handler should be a func")
	}

	if !FuncHaveExactlyTwoAttrib(handler) {
		panic("handler should receive a context and a command struct only")
	}

	if !FuncReturnsError(handler) {
		panic("handler should return error")
	}

	if !FuncAttributeMustBeContext(handler) {
		panic("handler first argument should have a context")
	}

	if !FuncAttributeMustBeStruct(handler) {
		panic("handler second argument must be a command struct")
	}
}
