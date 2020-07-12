package govalidator

import (
	"errors"
	"fmt"
	"reflect"
)

type constraintInvalid struct {
	k           reflect.Kind
	fi          *reflect.StructField
	requireFlag flagSet
	require     bool
}

func (c *constraintInvalid) reset() {
	c.k = reflect.Invalid
	c.requireFlag = setNo
	c.require = false
}

func (c *constraintInvalid) validate(value *reflect.Value, fix bool) error {
	if value.Kind() != reflect.Ptr {
		return nil
	}
	name := c.fi.Name
	// require
	if c.requireFlag == setYes && value.IsNil() {
		if fix && value.CanSet() {
			tye := value.Type().Elem()
			ns := reflect.New(tye)
			value.Set(ns)
		}
		return errors.New(fmt.Sprintf("`%s` is empty", name))
	}

	return nil
}

func (c *constraintInvalid) isSet() bool {
	return c.requireFlag == setYes
}

func describeInvalid(fi *reflect.StructField) (constraint, []error) {
	c := new(constraintInvalid)
	c.reset()
	es := make([]error, 0)
	//c.k =  fi.Type.Kind()
	c.k = getLastKind(fi.Type)
	c.fi = fi
	if req := fi.Tag.Get(flagReq); len(req) > 0 {
		c.requireFlag = setYes
		if req == "true" {
			c.require = true
		} else if req == "false" {
			c.require = false
		} else {
			c.require = false
			es = append(es, errors.New(fmt.Sprintf("`%s#req` value is invalid '%s', should be 'true' or 'false'", fi.Name, req)))
		}
	}

	return c, es
}
