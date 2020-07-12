package govalidator

import (
	"errors"
	"fmt"
	"math"
	"reflect"
)

const (
	floatThreshold = 1e-6
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

/*
* float compare
 */

func isZeroFloat64(f float64) bool {
	return math.Abs(f) < floatThreshold
}

func greaterFloat64Than(f1 float64, f2 float64) bool {
	diff := f1 - f2
	if isZeroFloat64(diff) {
		return false
	}
	// check f1 whether is zero
	if isZeroFloat64(f1) {
		if isZeroFloat64(f2) {
			return false
		}
		if f2 > floatThreshold {
			return false
		}
		return true
	}
	// check f2 whether is zero
	if isZeroFloat64(f2) {
		if f1 > floatThreshold {
			return true
		}
		return false
	}
	// compute precision, may be too big float value
	if isZeroFloat64(diff / f2) {
		return false
	}
	return diff > floatThreshold
}

func equalFloat64(f1 float64, f2 float64) bool {
	return !greaterFloat64Than(f1, f2) && !greaterFloat64Than(f2, f1)
}

func inSlice(ss []string, s string) bool {
	if len(ss) <= 0 {
		return false
	}
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}
