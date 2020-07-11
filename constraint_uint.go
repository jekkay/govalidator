package govalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type rangeUint struct {
	min uint64
	max uint64
}

type constraintUint struct {
	k  reflect.Kind
	fi *reflect.StructField

	minFlag     flagSet
	minUint     uint64
	maxFlag     flagSet
	maxUint     uint64
	defaultFlag flagSet
	defaultUint uint64
}

func (c *constraintUint) reset() {
	c.k = reflect.Invalid
	c.minFlag = set_no
	c.minUint = 0
	c.maxFlag = set_no
	c.maxUint = 0
	c.defaultFlag = set_no
	c.defaultUint = 0
}

func (c *constraintUint) validate(value *reflect.Value, fix bool) error {
	for value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil
		}
		e := value.Elem()
		value = &e
	}

	v := value.Uint()
	name := c.fi.Name
	if c.minFlag == set_yes && v < c.minUint {
		if fix && value.CanSet() {
			if v == 0 {
				value.SetUint(c.defaultUint)
			} else {
				value.SetUint(c.minUint)
			}
		}
		return errors.New(fmt.Sprintf("`%s` at least %d, current is %d", name, c.minUint, v))
	}

	if c.maxFlag == set_yes && v > c.maxUint {
		if fix && value.CanSet() {
			value.SetUint(c.maxUint)
		}
		return errors.New(fmt.Sprintf("`%s` at most %d, current is %d", name, c.maxUint, v))
	}

	return nil
}

func (c *constraintUint) isSet() bool {
	return c.minFlag == set_yes || c.maxFlag == set_yes || c.defaultFlag == set_yes
}

func describeUint(fi *reflect.StructField) (constraint, []error) {
	c := new(constraintUint)
	c.reset()
	es := make([]error, 0, 0)
	//c.k =  fi.Type.Kind()
	c.k = getLastKind(fi.Type)
	if _rangeUintMap[c.k] == nil {
		return nil, []error{errors.New(fmt.Sprintf("`%s` type is %v, required unsigned number type", fi.Name, c.k))}
	}
	c.fi = fi
	if minV := fi.Tag.Get(flagMin); len(minV) > 0 {
		if v, e := strconv.ParseUint(minV, 10, 64); e != nil {
			es = append(es, e)
		} else {
			c.minUint = v
			c.minFlag = set_yes
		}
	}
	if maxV := fi.Tag.Get(flagMax); len(maxV) > 0 {
		if v, e := strconv.ParseUint(maxV, 10, 64); e != nil {
			es = append(es, e)
		} else {
			c.maxUint = v
			c.maxFlag = set_yes
		}
	}
	if defaultV := fi.Tag.Get(flagDefault); len(defaultV) > 0 {
		if v, e := strconv.ParseUint(defaultV, 10, 64); e != nil {
			es = append(es, e)
		} else {
			c.defaultUint = v
			c.defaultFlag = set_yes
		}
	}
	if es2 := postCheckConstraintUint(c, fi); es2 != nil && len(es2) > 0 {
		es = append(es, es2...)
	}
	return c, es
}

func postCheckConstraintUint(c *constraintUint, fi *reflect.StructField) []error {
	name := fi.Name
	r := _rangeUintMap[c.k]
	if r == nil {
		return nil
	}
	es := make([]error, 0, 0)
	if c.minFlag == set_yes {
		if e := checkInRangeUint(name, flagMin, c.minUint, r); e != nil {
			es = append(es, e)
			c.minUint = r.min
		}
	}
	if c.maxFlag == set_yes {
		if e := checkInRangeUint(name, flagMax, c.maxUint, r); e != nil {
			es = append(es, e)
			c.maxUint = r.max
		}
	}

	if c.minFlag == set_yes && c.maxFlag == set_yes {
		if c.minUint > c.maxUint {
			es = append(es, errors.New(fmt.Sprintf("`%s` minimum value %d is greater than maximum value %d",
				name, c.minUint, c.maxUint)))
			c.maxUint = c.minUint
		}
	}

	if c.defaultFlag == set_yes {
		if e := checkInRangeUint(name, flagDefault, c.defaultUint, r); e != nil {
			es = append(es, e)
			c.defaultUint = 0
		}
		if c.minFlag == set_yes && c.defaultUint < c.minUint {
			es = append(es, errors.New(fmt.Sprintf("`%s#default` value is %d, should at least %d",
				name, c.defaultUint, c.minUint)))
			c.defaultUint = c.minUint
		} else if c.maxFlag == set_yes && c.defaultUint > c.maxUint {
			es = append(es, errors.New(fmt.Sprintf("`%s#default` value is %d, shold at most %d",
				name, c.defaultUint, c.maxUint)))
			c.defaultUint = c.maxUint
		}
	}
	if e := checkFirstLetter(fi, c); e != nil {
		es = append(es, e)
	}

	return es
}

func checkInRangeUint(name string, flag string, v uint64, r *rangeUint) error {
	if v < r.min {
		return errors.New(fmt.Sprintf("`%s#%s` value is %d, should at least %d", name, flag, v, r.min))
	}
	if v > r.max {
		return errors.New(fmt.Sprintf("`%s#%s` value is %d, should at most %d", name, flag, v, r.max))
	}
	return nil
}

var _rangeUintMap = make(map[reflect.Kind]*rangeUint)

func init() {
	_rangeUintMap[reflect.Uint8] = &rangeUint{min: 0, max: 255}
	_rangeUintMap[reflect.Uint16] = &rangeUint{min: 0, max: 65535}
	_rangeUintMap[reflect.Uint32] = &rangeUint{min: 0, max: 4294967295}
	_rangeUintMap[reflect.Uint] = &rangeUint{min: 0, max: 4294967295}
	_rangeUintMap[reflect.Uint64] = &rangeUint{min: 0, max: 18446744073709551615}
}
