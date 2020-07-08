package govalidator

import (
	"reflect"
	"strconv"
)

/**
* Number field range
 */
type FiledCategory int

const (
	UNKNOWN FiledCategory = iota
	SIGNED_NUMBER
	UNSIGNED_NUMBER
	STRING
)

type FiledTagInfo struct {
	Category FiledCategory
	Min      interface{}
	Max      interface{}
	Default  interface{}
}

func (fti *FiledTagInfo) IsValid() bool {
	v := false
	switch fti.Category {
	case SIGNED_NUMBER, UNSIGNED_NUMBER:
		v = fti.Max != nil || fti.Min == nil || fti.Default == nil
		//case STRING:
	}
	return v
}

func getRangeValue(fi *reflect.StructField) *FiledTagInfo {
	r := new(FiledTagInfo)
	if minV := fi.Tag.Get("min"); len(minV) > 0 {
		r.Min, _ = strconv.ParseInt(minV, 10, 64)
	}
	if maxV := fi.Tag.Get("max"); len(maxV) > 0 {
		r.Max, _ = strconv.ParseInt(maxV, 10, 64)
	}
	if defaultV := fi.Tag.Get("default"); len(defaultV) > 0 {
		r.Default, _ = strconv.ParseInt(defaultV, 10, 64)
	}
	r.Category = SIGNED_NUMBER
	return r
}

func getRangeUValue(fi *reflect.StructField) *FiledTagInfo {
	r := new(FiledTagInfo)
	if minV := fi.Tag.Get("min"); len(minV) > 0 {
		r.Min, _ = strconv.ParseUint(minV, 10, 64)
	}
	if maxV := fi.Tag.Get("max"); len(maxV) > 0 {
		r.Max, _ = strconv.ParseUint(maxV, 10, 64)
	}
	if defaultV := fi.Tag.Get("default"); len(defaultV) > 0 {
		r.Default, _ = strconv.ParseUint(defaultV, 10, 64)
	}
	r.Category = UNSIGNED_NUMBER
	return r
}
