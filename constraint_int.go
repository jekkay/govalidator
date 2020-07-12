package govalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type rangeInt struct {
	min int64
	max int64
}

type constraintInt struct {
	k  reflect.Kind
	fi *reflect.StructField

	minFlag     flagSet
	minInt      int64
	maxFlag     flagSet
	maxInt      int64
	defaultFlag flagSet
	defaultInt  int64
}

func (c *constraintInt) reset() {
	c.k = reflect.Invalid
	c.minFlag = setNo
	c.minInt = 0
	c.maxFlag = setNo
	c.maxInt = 0
	c.defaultFlag = setNo
	c.defaultInt = 0
}

func (c *constraintInt) validate(value *reflect.Value, fix bool) error {
	for value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil
		}
		e := value.Elem()
		value = &e
	}

	v := value.Int()
	name := c.fi.Name
	if c.minFlag == setYes && v < c.minInt {
		if fix && value.CanSet() {
			if v == 0 {
				value.SetInt(c.defaultInt)
			} else {
				value.SetInt(c.minInt)
			}
		}
		return errors.New(fmt.Sprintf("`%s` at least %d, current is %d", name, c.minInt, v))
	}

	if c.maxFlag == setYes && v > c.maxInt {
		if fix && value.CanSet() {
			value.SetInt(c.maxInt)
		}
		return errors.New(fmt.Sprintf("`%s` at most %d, current is %d", name, c.maxInt, v))
	}

	return nil
}

func (c *constraintInt) isSet() bool {
	return c.minFlag == setYes || c.maxFlag == setYes || c.defaultFlag == setYes
}

func describeInt(fi *reflect.StructField) (constraint, []error) {
	c := new(constraintInt)
	c.reset()
	es := make([]error, 0)
	//c.k =  fi.Type.Kind()
	c.k = getLastKind(fi.Type)
	if _rangeIntMap[c.k] == nil {
		return nil, []error{errors.New(fmt.Sprintf("`%s` type is %v, required signed number type", fi.Name, c.k))}
	}
	c.fi = fi
	if minV := fi.Tag.Get(flagMin); len(minV) > 0 {
		if v, e := strconv.ParseInt(minV, 10, 64); e != nil {
			es = append(es, e)
		} else {
			c.minInt = v
			c.minFlag = setYes
		}
	}
	if maxV := fi.Tag.Get(flagMax); len(maxV) > 0 {
		if v, e := strconv.ParseInt(maxV, 10, 64); e != nil {
			es = append(es, e)
		} else {
			c.maxInt = v
			c.maxFlag = setYes
		}
	}
	if defaultV := fi.Tag.Get(flagDefault); len(defaultV) > 0 {
		if v, e := strconv.ParseInt(defaultV, 10, 64); e != nil {
			es = append(es, e)
		} else {
			c.defaultInt = v
			c.defaultFlag = setYes
		}
	}
	if es2 := postCheckConstraintInt(c, fi); es2 != nil && len(es2) > 0 {
		es = append(es, es2...)
	}
	return c, es
}

func postCheckConstraintInt(c *constraintInt, fi *reflect.StructField) []error {
	name := fi.Name
	r := _rangeIntMap[c.k]
	if r == nil {
		return nil
	}
	es := make([]error, 0, 0)
	if c.minFlag == setYes {
		if e := checkInRangeInt(name, flagMin, c.minInt, r); e != nil {
			es = append(es, e)
			c.minInt = r.min
		}
	}
	if c.maxFlag == setYes {
		if e := checkInRangeInt(name, flagMax, c.maxInt, r); e != nil {
			es = append(es, e)
			c.maxInt = r.max
		}
	}

	if c.minFlag == setYes && c.maxFlag == setYes {
		if c.minInt > c.maxInt {
			es = append(es, errors.New(fmt.Sprintf("`%s` minimum value %d is greater than maximum value %d",
				name, c.minInt, c.maxInt)))
			c.maxInt = c.minInt
		}
	}

	if c.defaultFlag == setYes {
		if e := checkInRangeInt(name, flagDefault, c.defaultInt, r); e != nil {
			es = append(es, e)
			c.defaultInt = 0
		}
		if c.minFlag == setYes && c.defaultInt < c.minInt {
			es = append(es, errors.New(fmt.Sprintf("`%s#default` value is %d, should at least %d",
				name, c.defaultInt, c.minInt)))
			c.defaultInt = c.minInt
		} else if c.maxFlag == setYes && c.defaultInt > c.maxInt {
			es = append(es, errors.New(fmt.Sprintf("`%s#default` value is %d, shold at most %d",
				name, c.defaultInt, c.maxInt)))
			c.defaultInt = c.maxInt
		}
	}
	if e := checkFirstLetter(fi, c); e != nil {
		es = append(es, e)
	}

	return es
}

func checkInRangeInt(name string, flag string, v int64, r *rangeInt) error {
	if v < r.min {
		return errors.New(fmt.Sprintf("`%s#%s` value is %d, should at least %d", name, flag, v, r.min))
	}
	if v > r.max {
		return errors.New(fmt.Sprintf("`%s#%s` value is %d, should at most %d", name, flag, v, r.max))
	}
	return nil
}

var _rangeIntMap = make(map[reflect.Kind]*rangeInt)

func init() {
	_rangeIntMap[reflect.Int8] = &rangeInt{min: -128, max: 127}
	_rangeIntMap[reflect.Int16] = &rangeInt{min: -32768, max: 32767}
	_rangeIntMap[reflect.Int32] = &rangeInt{min: -2147483648, max: 2147483647}
	_rangeIntMap[reflect.Int] = &rangeInt{min: -2147483648, max: 2147483647}
	_rangeIntMap[reflect.Int64] = &rangeInt{min: -9223372036854775808, max: 9223372036854775807}
}
