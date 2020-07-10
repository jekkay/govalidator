package govalidator

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFiledTag_describe(t *testing.T) {
	type sa struct {
		Age int8 `min:"1000" max:"200" default:"10"`
		Score *int16 `min:"-100000000" max:"200" default:"999"`
	}

	a := &sa{Age:1}
	v := reflect.ValueOf(a).Elem()
	for i:= 0; i < v.NumField(); i ++ {
		fi := v.Type().Field(i)
		if ds, er := describeInt(&fi); er != nil {
			fmt.Println(er)
			fmt.Println(ds)
		} else {
			t.Errorf("fail to detect describe errors")
		}
	}
}
