package govalidator

import (
	"errors"
	"fmt"
	"reflect"
)

func getLastKind(t reflect.Type) reflect.Kind {
	if k := t.Kind(); k != reflect.Ptr {
		return k
	}
	return getLastKind(t.Elem())
}

func checkFirstLetter(fi *reflect.StructField, c constraint) error {
	if !c.isSet() {
		return nil
	}
	if name := fi.Name; name[0] < 'A' || name[0] > 'Z' {
		return errors.New(fmt.Sprintf("`%s` first letter should be upper case", name))
	}
	return nil
}
