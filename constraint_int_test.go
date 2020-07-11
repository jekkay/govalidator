package govalidator

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFiledTag_describe(t *testing.T) {
	type sa struct {
		Age   int8   `min:"1000" max:"20" default:"10"`
		Score *int16 `min:"1" max:"200000" default:"99"`
	}

	a := &sa{Age: 50, Score: new(int16)}
	v := reflect.ValueOf(a).Elem()
	for i := 0; i < v.NumField(); i++ {
		fi := v.Type().Field(i)
		ds, er := describeInt(&fi)
		if er != nil {
			fmt.Println(er)
			fmt.Println(ds)
		} else {
			t.Fatal("fail to detect describe errors")
		}

		if ds.isSet() {
			val := v.Field(i)
			if es := ds.validate(&val, true); es != nil {
				fmt.Println(es)
			}

		}
	}
	fmt.Printf("Age: %d, Score: %d\n", a.Age, *a.Score)
	if a.Age != 20 {
		t.Errorf("fail to auto fix filed Age")
	}
	if *a.Score != 99 {
		t.Errorf("fail to auto fix filed Score")
	}
}

func TestDescribe_invalidType(t *testing.T) {
	type sa struct {
		Age uint16 `min:"10" max:"20" default:"15"`
	}
	a := &sa{Age: 50}
	v := reflect.ValueOf(a).Elem()
	fi := v.Type().Field(0)
	if _, er := describeInt(&fi); er != nil {
		fmt.Printf("%v", er)
	} else {
		t.Errorf("fail to find type unmatch error")
	}
}

func TestDescribe_invalidVariable(t *testing.T) {
	type sa struct {
		age int16 `min:"10" max:"20" default:"15"`
	}
	a := &sa{age: 50}
	v := reflect.ValueOf(a).Elem()
	fi := v.Type().Field(0)
	if _, er := describeInt(&fi); er != nil {
		fmt.Printf("%v", er)
	} else {
		t.Errorf("fail to find variable error")
	}
}
