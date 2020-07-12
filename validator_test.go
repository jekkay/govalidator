package govalidator

import (
	"encoding/json"
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

	if a.Age != 50 || a.Year != 0 || *a.Score != 40 || a.Location != "CN" || *a.Name != "jekkay" {
		t.Error("fail to adjust values")
	}
}

func TestValidator_Validates(t *testing.T) {
	type A struct {
		Age  int8  `min:"-100" max:"3000" default:"10"`
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

func TestValidObject2(t *testing.T) {
	type Range struct {
		A int32   `json:"a" min:"10" max:"100" default:"50"`
		B int32   `json:"b" min:"20" max:"90" default:"80"`
		C *uint64 `json:"c" min:"30" max:"90" req:"true" default:"60"`
	}

	r := new(Range)
	r.A = 120
	r.B = 130

	if e := ValidObject(r, false); e != nil {
		fmt.Println(e)
	}
	ValidObject(r, true)
	bs, _ := json.MarshalIndent(r, "", "  ")
	fmt.Println(string(bs))
}
