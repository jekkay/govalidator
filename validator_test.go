package govalidator

import (
	"fmt"
	"testing"
)

func TestValidator_Validate(t *testing.T) {
	type A struct {
		Age      int8    `min:"10" max:"100" default:"20"`
		Year     int16   `min:"-100" max:"9999" default:"30"`
		Score    *int16  `min:"20" max:"99" default:"40" req:"true"`
		Location string  `min:"1" max:"5" default:"CN" req:"true"`
		Name     *string `min:"2" max:"10" default:"jekkay" req:"true"`
	}

	a := &A{Age: 50}
	if e := ValidObject(a, true); e != nil {
		fmt.Println(e)
	}

	if a.Age != 50 || a.Year != 0 || *a.Score != 40 ||a.Location != "CN" || *a.Name != "jekkay"{
		t.Error("fail to adjust values")
	}
}

func TestValidator_Validates(t *testing.T) {
	type A struct {
		Age int8 `min:"-100" max:"3000" default:"10"`
		Year int16 `min:"0" max:"100000" default:"20"`
	}
	a := &A{}
	if e := Validate(a); e != nil {
		fmt.Println(e)
	} else {
		t.Error("fail to find errors")
	}

	a2 := &A{}
	if es := Validates(a2); len(es) > 0 {
		fmt.Println(es)
	} else {
		t.Error("Validates: fail to find errors")
	}
}
