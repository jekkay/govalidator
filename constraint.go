package govalidator

import "reflect"

/**
*  constraint interface for different filed type
 */

type constraint interface {
	validate(value *reflect.Value, fix bool) error
	isSet() bool
}

type tagDescriber func(fi *reflect.StructField) (constraint, []error)

