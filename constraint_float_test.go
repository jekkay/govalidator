package govalidator

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFiledTagInfo_range(t *testing.T) {
	type sa struct {
		Div float32 `json:"age" min:"-5e70" max:"2.5e20" default:"-1.8e20"`
	}
	a := &sa{Div: 2.6e22}
	v := reflect.ValueOf(a).Elem()
	fi := v.Type().Field(0)
	if d, er := describeFloat(&fi); len(er) <= 0 {
		t.Error("fail to find minimum value error")
	} else {
		fmt.Println(er)
		val := v.Field(0)
		if e := d.validate(&val, true); e != nil {
			fmt.Sprint(e)
		}
	}

	if !equalFloat64(float64(a.Div), 2.5e20) {
		t.Error("fail to adjust value")
	}
}
