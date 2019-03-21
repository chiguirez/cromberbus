package typer

import "reflect"

func Identify(i interface{}) string {
	rType := reflect.TypeOf(i)
	if rType.Kind() == reflect.Ptr {
		return rType.Elem().Name()
	}

	return rType.Name()
}
