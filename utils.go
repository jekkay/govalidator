package govalidator

import "reflect"

func getLastKind(t reflect.Type) reflect.Kind {
	if k := t.Kind(); k != reflect.Ptr {
		return k
	}
	return getLastKind(t.Elem())
}
