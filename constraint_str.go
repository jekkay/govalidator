package govalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type constraintString struct {
	k  reflect.Kind
	fi *reflect.StructField

	minFlag      flagSet
	minLen       int
	maxFlag      flagSet
	maxLen       int
	defaultFlag  flagSet
	defaultStr   string
	requireFlag  flagSet
	require      bool
	inFlag       flagSet
	in           []string
	RegExFlag    flagSet
	RegEx        string
	RegExCompile *regexp.Regexp
}

func (c *constraintString) reset() {
	c.k = reflect.Invalid
	c.minFlag = setNo
	c.minLen = 0
	c.maxFlag = setNo
	c.maxLen = 0
	c.defaultFlag = setNo
	c.defaultStr = ""
	c.requireFlag = setNo
	c.require = false
	c.inFlag = setNo
	c.in = nil
	c.RegExFlag = setNo
	c.RegEx = ""
	c.RegExCompile = nil
}

func (c *constraintString) validate(value *reflect.Value, fix bool) error {
	for value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil
		}
		e := value.Elem()
		value = &e
	}

	v := value.String()
	l := len(v)
	name := c.fi.Name
	// check empty string
	// require
	if c.requireFlag == setYes && l == 0 {
		if fix && value.CanSet() {
			value.SetString(c.defaultStr)
		}
		return errors.New(fmt.Sprintf("`%s` is empty", name))
	} else if c.requireFlag == setNo && l == 0 {
		// empty string
		return nil
	}
	// string length check
	if c.minFlag == setYes && l < c.minLen {
		if fix && value.CanSet() {
			value.SetString(c.defaultStr)
		}
		return errors.New(fmt.Sprintf("`%s` at least length %d, current length is %d, %s", name, c.minLen, l, v))
	}

	if c.maxFlag == setYes && l > c.maxLen {
		if fix && value.CanSet() {
			value.SetString(c.defaultStr)
		}
		return errors.New(fmt.Sprintf("`%s` at most length %d, current length is %d, %s", name, c.maxLen, l, v))
	}
	// in options
	if c.inFlag == setYes && !inSlice(c.in, v) {
		if fix && value.CanSet() {
			value.SetString(c.defaultStr)
		}
		return errors.New(fmt.Sprintf("`%s` value '%s' is not valid, should be in options:%v", name, v, c.in))
	}
	// regex
	if c.RegExFlag == setYes && len(c.RegExCompile.FindString(v)) == 0 {
		if fix && value.CanSet() {
			value.SetString(c.defaultStr)
		}
		return errors.New(fmt.Sprintf("`%s` value '%s' doesn't match '%s'", name, v, c.RegEx))
	}

	return nil
}

func (c *constraintString) isSet() bool {
	return c.minFlag == setYes || c.maxFlag == setYes || c.defaultFlag == setYes ||
		c.requireFlag == setYes || c.RegExFlag == setYes || c.inFlag == setYes
}

func describeString(fi *reflect.StructField) (constraint, []error) {
	c := new(constraintString)
	c.reset()
	es := make([]error, 0, 0)
	//c.k =  fi.Type.Kind()
	c.k = getLastKind(fi.Type)
	c.fi = fi
	if minV := fi.Tag.Get(flagMin); len(minV) > 0 {
		if v, e := strconv.ParseInt(minV, 10, 32); e != nil {
			es = append(es, e)
		} else {
			c.minLen = int(v)
			c.minFlag = setYes
		}
	}
	if maxV := fi.Tag.Get(flagMax); len(maxV) > 0 {
		if v, e := strconv.ParseInt(maxV, 10, 32); e != nil {
			es = append(es, e)
		} else {
			c.maxLen = int(v)
			c.maxFlag = setYes
		}
	}
	if defaultV := fi.Tag.Get(flagDefault); len(defaultV) > 0 {
		c.defaultStr = defaultV
		c.defaultFlag = setYes
	}
	if in := fi.Tag.Get(flagIn); len(in) > 0 {
		c.inFlag = setYes
		c.in = stripEmptyString(splitInOptions(in))
	}
	if req := fi.Tag.Get(flagReq); len(req) > 0 {
		c.requireFlag = setYes
		if req == "true" {
			c.require = true
		} else if req == "false" {
			c.require = false
		} else {
			c.require = false
			es = append(es, errors.New(fmt.Sprintf("`%s#req` value is invalid '%s', should be 'true' or 'false'", fi.Name, req)))
		}
	}
	if regex := fi.Tag.Get(flagRegEx); len(regex) > 0 {
		c.RegEx = regex
		c.RegExFlag = setYes
	}
	if es2 := postCheckConstraintString(c, fi); es2 != nil && len(es2) > 0 {
		es = append(es, es2...)
	}
	return c, es
}

func splitInOptions(in string) []string {
	if len(in) <= 0 {
		return nil
	}
	// check the first character whether is '#'
	if in[0] == strSep1[0] {
		return strings.Split(in, strSep1)
	}
	return strings.Split(in, strSep2)
}

func stripEmptyString(ss []string) []string {
	if len(ss) <= 0 {
		return nil
	}
	r := make([]string, 0)
	for _, s := range ss {
		if len(s) > 0 {
			r = append(r, s)
		}
	}
	return r
}

func postCheckConstraintString(c *constraintString, fi *reflect.StructField) []error {
	if !c.isSet() {
		return nil
	}
	name := fi.Name
	es := make([]error, 0, 0)

	if c.minFlag == setYes && c.maxFlag == setYes {
		if c.minLen > c.maxLen {
			es = append(es, errors.New(fmt.Sprintf("`%s` minimum length %d is greater than maximum length %d",
				name, c.minLen, c.maxLen)))
			c.maxLen = c.minLen
		}
	}
	if e := checkFirstLetter(fi, c); e != nil {
		es = append(es, e)
	}
	// regular expression compile
	if c.RegExFlag == setYes {
		if r, e := regexp.Compile(c.RegEx); e != nil {
			es = append(es, e)
			c.RegExFlag = setNo
		} else {
			c.RegExCompile = r
		}
	}

	return es
}
