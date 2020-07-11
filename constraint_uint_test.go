package govalidator

import (
	"reflect"
	"testing"
)

func TestFiledTagInfo_IsValid(t *testing.T) {
	type sa struct {
		age uint8 `json:"age"`
	}
	a := &sa{age: 50}
	v := reflect.ValueOf(a).Elem()
	fi := v.Type().Field(0)
	if _, er := describeUint(&fi); len(er) > 0 {
		t.Error(er)
	}
}

func TestFiledTagInfo_IsValid2(t *testing.T) {
	type sa struct {
		age uint8 `json:"age" min:"100"`
	}
	a := &sa{age: 50}
	v := reflect.ValueOf(a).Elem()
	fi := v.Type().Field(0)
	if _, er := describeUint(&fi); len(er) > 0 {
		t.Error(er)
	}
}

func TestFiledTagInfo_adjust(t *testing.T) {
	type sa struct {
		Age uint8 `json:"age" min:"100"`
	}
	a := &sa{Age: 50}
	v := reflect.ValueOf(a).Elem()
	fi := v.Type().Field(0)
	if c, er := describeUint(&fi); len(er) > 0 {
		t.Error(er)
	} else if c != nil {
		va := v.Field(0)
		if e := c.validate(&va, true); e == nil {
			t.Error("fail to find value error")
		}
	}
	if a.Age != 100 {
		t.Error("fail to change uint value to min")
	}
}

func TestFiledTagInfo_adjust2(t *testing.T) {
	type sa struct {
		Age uint8 `json:"age" min:"100" max:"200" default:"150"`
	}
	a := &sa{}
	v := reflect.ValueOf(a).Elem()
	fi := v.Type().Field(0)
	if c, er := describeUint(&fi); len(er) > 0 {
		t.Error(er)
	} else if c != nil {
		va := v.Field(0)
		if e := c.validate(&va, true); e == nil {
			t.Error("fail to find value error")
		}
	}
	if a.Age != 150 {
		t.Error("fail to auto fix uint value to default")
	}
}
