package govalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type rangeFloat struct {
	min float64
	max float64
}

type constraintFloat struct {
	k  reflect.Kind
	fi *reflect.StructField

	minFlag      flagSet
	minFloat     float64
	maxFlag      flagSet
	maxFloat     float64
	defaultFlag  flagSet
	defaultFloat float64
}

func (c *constraintFloat) reset() {
	c.k = reflect.Invalid
	c.minFlag = setNo
	c.minFloat = 0
	c.maxFlag = setNo
	c.maxFloat = 0
	c.defaultFlag = setNo
	c.defaultFloat = 0
}

func (c *constraintFloat) validate(value *reflect.Value, fix bool) error {
	for value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil
		}
		e := value.Elem()
		value = &e
	}

	v := value.Float()
	name := c.fi.Name
	if c.minFlag == setYes && greaterFloat64Than(c.minFloat, v) {
		if fix && value.CanSet() {
			if isZeroFloat64(v) {
				value.SetFloat(c.defaultFloat)
			} else {
				value.SetFloat(c.minFloat)
			}
		}
		return errors.New(fmt.Sprintf("`%s` at leas %.5f, current is %.5f", name, c.minFloat, v))
	}

	if c.maxFlag == setYes && greaterFloat64Than(v, c.maxFloat) {
		if fix && value.CanSet() {
			value.SetFloat(c.maxFloat)
		}
		return errors.New(fmt.Sprintf("`%s` at most %.5f, current is %.5f", name, c.maxFloat, v))
	}

	return nil
}

func (c *constraintFloat) isSet() bool {
	return c.minFlag == setYes || c.maxFlag == setYes || c.defaultFlag == setYes
}

func describeFloat(fi *reflect.StructField) (constraint, []error) {
	c := new(constraintFloat)
	c.reset()
	es := make([]error, 0, 0)
	//c.k =  fi.Type.Kind()
	c.k = getLastKind(fi.Type)
	if _rangeFloatMap[c.k] == nil {
		return nil, []error{errors.New(fmt.Sprintf("`%s` type is %v, required float type", fi.Name, c.k))}
	}
	c.fi = fi
	if minV := fi.Tag.Get(flagMin); len(minV) > 0 {
		if v, e := strconv.ParseFloat(minV, 64); e != nil {
			es = append(es, e)
		} else {
			c.minFloat = v
			c.minFlag = setYes
		}
	}
	if maxV := fi.Tag.Get(flagMax); len(maxV) > 0 {
		if v, e := strconv.ParseFloat(maxV, 64); e != nil {
			es = append(es, e)
		} else {
			c.maxFloat = v
			c.maxFlag = setYes
		}
	}
	if defaultV := fi.Tag.Get(flagDefault); len(defaultV) > 0 {
		if v, e := strconv.ParseFloat(defaultV, 64); e != nil {
			es = append(es, e)
		} else {
			c.defaultFloat = v
			c.defaultFlag = setYes
		}
	}
	if es2 := postCheckConstraintFloat(c, fi); es2 != nil && len(es2) > 0 {
		es = append(es, es2...)
	}
	return c, es
}

func postCheckConstraintFloat(c *constraintFloat, fi *reflect.StructField) []error {
	name := fi.Name
	r := _rangeFloatMap[c.k]
	if r == nil {
		return nil
	}
	es := make([]error, 0, 0)
	if c.minFlag == setYes {
		if e := checkInRangeFloat(name, flagMin, c.minFloat, r); e != nil {
			es = append(es, e)
			c.minFloat = r.min
		}
	}
	if c.maxFlag == setYes {
		if e := checkInRangeFloat(name, flagMax, c.maxFloat, r); e != nil {
			es = append(es, e)
			c.maxFloat = r.max
		}
	}

	if c.minFlag == setYes && c.maxFlag == setYes {
		if greaterFloat64Than(c.minFloat, c.maxFloat) {
			es = append(es, errors.New(fmt.Sprintf("`%s` minimum value %.5f is greater than maximum value %.5f",
				name, c.minFloat, c.maxFloat)))
			c.maxFloat = c.minFloat
		}
	}

	if c.defaultFlag == setYes {
		if e := checkInRangeFloat(name, flagDefault, c.defaultFloat, r); e != nil {
			es = append(es, e)
			c.defaultFloat = 0
		}
		if c.minFlag == setYes && greaterFloat64Than(c.minFloat, c.defaultFloat) {
			es = append(es, errors.New(fmt.Sprintf("`%s#default` value is %.5f, should at least %.5f",
				name, c.defaultFloat, c.minFloat)))
			c.defaultFloat = c.minFloat
		} else if c.maxFlag == setYes && greaterFloat64Than(c.defaultFloat, c.maxFloat) {
			es = append(es, errors.New(fmt.Sprintf("`%s#default` value is %.5f, shold at most %.5f",
				name, c.defaultFloat, c.maxFloat)))
			c.defaultFloat = c.maxFloat
		}
	}
	if e := checkFirstLetter(fi, c); e != nil {
		es = append(es, e)
	}

	return es
}

func checkInRangeFloat(name string, flag string, v float64, r *rangeFloat) error {
	if greaterFloat64Than(r.min, v) {
		return errors.New(fmt.Sprintf("`%s#%s` value is %.5f, should at least %.5f", name, flag, v, r.min))
	}
	if greaterFloat64Than(v, r.max) {
		return errors.New(fmt.Sprintf("`%s#%s` value is %.5f, should at most %.5f", name, flag, v, r.max))
	}
	return nil
}

var _rangeFloatMap = make(map[reflect.Kind]*rangeFloat)

func init() {
	_rangeFloatMap[reflect.Float32] = &rangeFloat{min: -3.4e38, max: 3.4e38}
	_rangeFloatMap[reflect.Float64] = &rangeFloat{min: -1.7e308, max: 1.7e308}
}
