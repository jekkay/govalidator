package govalidator

import (
	"errors"
	"reflect"
)

/**
* global tag describer map
 */

var constraintMap map[reflect.Kind]tagDescriber

type Validator struct {
}

/**
* create validator instance
 */

func (r *Validator) Validate(obj interface{}) error {
	if es := r.validInner(obj, false, true); es != nil && len(es) > 0 {
		return es[0]
	}
	return nil
}

func (r *Validator) Validates(obj interface{}) []error {
	return r.validInner(obj, false, false)
}

func (r *Validator) ValidObject(obj interface{}, fix bool) []error {
	return r.validInner(obj, fix, false)
}

func (r *Validator) validInner(obj interface{}, fix bool, quick bool) []error {
	t := reflect.TypeOf(obj)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return []error{errors.New("object must be struct pointer")}
	}
	v := reflect.ValueOf(obj).Elem()
	return r.doValidObject(&v, fix, quick)
}

func (r *Validator) doValidObject(v *reflect.Value, fix bool, quick bool) []error {
	if v.Kind() == reflect.Invalid {
		return nil
	}
	es := make([]error, 0)
	// only process struct
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			fi := v.Type().Field(i)
			val := v.Field(i)
			if errs := r.checkFieldValue(&fi, &val, fix, quick); errs != nil && len(errs) > 0 {
				es = append(es, errs...)
				if quick {
					return errs
				}
			}
		}
	}
	return es
}

func (r *Validator) checkFieldValue(fi *reflect.StructField, value *reflect.Value, fix bool, quick bool) []error {
	if value.Kind() == reflect.Invalid {
		return nil
	}
	switch value.Kind() {
	case reflect.Ptr:
		es := make([]error, 0)
		v2 := value.Elem()
		if v2.Kind() == reflect.Invalid {
			c, e1 := describeInvalid(fi)
			if len(e1) > 0 {
				es = append(es, e1...)
			}
			if c != nil && c.isSet() {
				if e2 := c.validate(value, fix); e2 != nil {
					es = append(es, e2)
				}
				// re-try get new instance
				v2 = value.Elem()
			} else {
				return nil
			}
		}
		if v2.Kind() == reflect.Struct {
			if e3 := r.doValidObject(&v2, fix, quick); len(e3) > 0 {
				es = append(es, e3...)
			}
		} else if e3 := r.checkFieldValue(fi, &v2, fix, quick); len(e3) > 0 {
			es = append(es, e3...)
		}
		return es
	case reflect.Struct:
		return r.doValidObject(value, fix, quick)
	default:
		es := make([]error, 0)
		if f := constraintMap[value.Kind()]; f != nil {
			c, e1 := f(fi)
			if e1 != nil && len(e1) > 0 {
				es = append(es, e1...)
			} else if c != nil && c.isSet() {
				if e2 := r.validateFieldValue(c, value, fix); e2 != nil && len(e2) > 0 {
					es = append(es, e2...)
				}
			}
		} else {
			// not support yet
		}
		if len(es) > 0 {
			return es
		}
	}
	return nil
}

/**
* validate field, auto adjust field value if `fix` is true
 */
func (r *Validator) validateFieldValue(c constraint, value *reflect.Value, fix bool) []error {
	if e := c.validate(value, fix); e != nil {
		return []error{e}
	}
	return nil
}

/**
* register constraint data model when startup
 */
func init() {
	constraintMap = make(map[reflect.Kind]tagDescriber)
	constraintMap[reflect.Int] = describeInt
	constraintMap[reflect.Int8] = describeInt
	constraintMap[reflect.Int16] = describeInt
	constraintMap[reflect.Int32] = describeInt
	constraintMap[reflect.Int64] = describeInt

	constraintMap[reflect.Uint] = describeUint
	constraintMap[reflect.Uint8] = describeUint
	constraintMap[reflect.Uint16] = describeUint
	constraintMap[reflect.Uint32] = describeUint
	constraintMap[reflect.Uint64] = describeUint

	constraintMap[reflect.Float32] = describeFloat
	constraintMap[reflect.Float64] = describeFloat

	constraintMap[reflect.String] = describeString
}
