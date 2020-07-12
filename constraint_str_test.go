package govalidator

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestFiledTagInfo_checkMin(t *testing.T) {
	type sa struct {
		Msg string `json:"msg" min:"2" default:"hello"`
	}
	a := &sa{Msg: "a"}
	v := reflect.ValueOf(a).Elem()
	fi := v.Type().Field(0)
	if d, er := describeString(&fi); len(er) > 0 {
		t.Error(er)
	} else {
		val := v.Field(0)
		if e := d.validate(&val, true); e != nil {
			fmt.Print(e)
		} else {
			t.Error("fail to adjust string value")
		}
	}
	if a.Msg != "hello" {
		t.Error("adjust value error")
	}
}

func TestFiledTagInfo_InTag1(t *testing.T) {
	type sa struct {
		Msg string `json:"msg" min:"2" in:"hello,world,jekkay" default:"jekkay"`
	}
	a := &sa{}
	v := reflect.ValueOf(a).Elem()
	fi := v.Type().Field(0)
	if d, er := describeString(&fi); len(er) > 0 {
		t.Error(er)
	} else {
		val := v.Field(0)
		if e := d.validate(&val, true); e != nil {
			t.Error(e)
		}
	}
}

func TestFiledTagInfo_InTag2(t *testing.T) {
	type sa struct {
		Msg string `json:"msg" min:"2" in:"hello,world,jekkay" default:"jekkay" req:"true"`
	}
	a := &sa{}
	v := reflect.ValueOf(a).Elem()
	fi := v.Type().Field(0)
	if d, er := describeString(&fi); len(er) > 0 {
		t.Error(er)
	} else {
		val := v.Field(0)
		if e := d.validate(&val, true); e != nil {
			fmt.Println(e)
		} else {
			t.Error("fail to adjust string to default value")
		}
	}
	if a.Msg != "jekkay" {
		t.Error("fail to adjust string to default value 'jekkay'")
	}
}

func TestFiledTagInfo_InTag3(t *testing.T) {
	type sa struct {
		Msg string `json:"msg" min:"2" in:"hello,world,jekkay" default:"jekkay" req:"true"`
	}
	a := &sa{Msg: "good"}
	v := reflect.ValueOf(a).Elem()
	fi := v.Type().Field(0)
	if d, er := describeString(&fi); len(er) > 0 {
		t.Error(er)
	} else {
		val := v.Field(0)
		if e := d.validate(&val, true); e != nil {
			fmt.Println(e)
		} else {
			t.Error("fail to adjust string to default value")
		}
	}
	if a.Msg != "jekkay" {
		t.Error("fail to adjust string to default value 'jekkay'")
	}
}

func TestFiledTagInfo_InTag4(t *testing.T) {
	type sa struct {
		Msg string `json:"msg" min:"2" regex:"^[a-z]+\\d+$" default:"jekkay12" req:"true"`
	}
	a := &sa{Msg: "good12"}
	v := reflect.ValueOf(a).Elem()
	fi := v.Type().Field(0)
	if d, er := describeString(&fi); len(er) > 0 {
		t.Error(er)
	} else {
		val := v.Field(0)
		if e := d.validate(&val, true); e != nil {
			t.Error(e)
		}
	}
}

func TestFiledTagInfo_InTag5(t *testing.T) {
	type sa struct {
		Msg string `json:"msg" min:"2" regex:"^[a-z]+\\d+$" default:"jekkay12" req:"true"`
	}
	a := &sa{Msg: "good"}
	v := reflect.ValueOf(a).Elem()
	fi := v.Type().Field(0)
	if d, er := describeString(&fi); len(er) > 0 {
		t.Error(er)
	} else {
		val := v.Field(0)
		if e := d.validate(&val, true); e != nil {
			fmt.Println(e)
		} else {
			t.Error(errors.New("fail to find mismatch regex error"))
		}
	}

	if a.Msg != "jekkay12" {
		t.Error("fail to adjust value")
	}
}

