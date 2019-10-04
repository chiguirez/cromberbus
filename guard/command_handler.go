package guard

import "reflect"

func FuncHaveMoreThanOneAttrib(handler interface{}) bool {
	return reflect.TypeOf(handler).NumIn() > 1
}

func FuncAttributeMustBeStruct(handler interface{}) bool {
	return reflect.TypeOf(handler).In(0).Kind() == reflect.Struct
}

func IsAFunc(handler interface{}) bool {
	return reflect.TypeOf(handler).Kind() == reflect.Func
}

func CommandHandler(handler interface{}) {
	if !IsAFunc(handler) {
		panic("handler should be a func")
	}
	if FuncHaveMoreThanOneAttrib(handler) {
		panic("handler should only receive one input")
	}
	if !FuncAttributeMustBeStruct(handler) {
		panic("handler input must be a command struct")
	}
}
