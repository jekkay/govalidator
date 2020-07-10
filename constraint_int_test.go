package govalidator

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFiledTag_describe(t *testing.T) {
	type sa struct {
		Age int8 `min:"1000" max:"20" default:"10"`
		//Score *int16 `min:"-100000000" max:"200" default:"999"`
	}

	a := &sa{Age:50}
	v := reflect.ValueOf(a).Elem()
	for i:= 0; i < v.NumField(); i ++ {
		fi := v.Type().Field(i)
		ds, er := describeInt(&fi)
		if er != nil {
			fmt.Println(er)
			fmt.Println(ds)
		} else {
			t.Fatal("fail to detect describe errors")
		}

		if  ds.isSet() {
			val := v.Field(i)
			if es := ds.validate(&val, true); es != nil {
				fmt.Println(es)
			}
			if a.Age != 20 {
				t.Errorf("fail to auto fix filed Age")
			}
		}
	}
}
