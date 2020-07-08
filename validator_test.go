package govalidator

import (
	"encoding/json"
	"fmt"
	"testing"
)

type Range struct {
	A int32   `json:"a" min:"10" max:"100" default:"50"`
	B int32   `json:"b" min:"20" max:"90" default:"80"`
	C *uint64 `json:"c" min:"30" max:"90" default:"60"`
}

type Student struct {
	Age   uint16 `json:"age" min:"0" max:"100" default:"20"`
	Year  int32  `json:"year" min:"1000" max:"9999" default:"2018"`
	Grade int64  `json:"grade" min:"1" max:"9"`
	gg    int64  `json:"gg" min:"1" max:"9"`
	R1    Range  `json:"r1"`
	R2    *Range `json:"r2"`
}

func TestValidObject2(t *testing.T) {
	r := new(Range)
	r.A = 120
	r.B = 130
	r.C = new(uint64)
	*r.C = 0

	if e := ValidObject(r, false); e != nil {
		fmt.Println(e)
	}
	ValidObject(r, true)
	bs, _ := json.MarshalIndent(r, "", "  ")
	fmt.Println(string(bs))
}

func TestValidObject(t *testing.T) {
	s := new(Student)
	s.Age = 200
	s.Year = 20
	s.Grade = -1
	s.gg = 0
	s.R2 = new(Range)
	s.R2.A = 300
	s.R2.B = 1
	s.R2.C = new(uint64)
	*s.R2.C = 400
	if bs, _ := json.MarshalIndent(s, "", "  "); len(bs) > 0 {
		fmt.Println(string(bs))
	}
	if e := ValidObject(s, true); e != nil {
		fmt.Println(e)
	}
	fmt.Println("After fix filed value")
	if bs, _ := json.MarshalIndent(s, "", "  "); len(bs) > 0 {
		fmt.Println(string(bs))
	}

}
