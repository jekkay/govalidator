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
	c.minFlag = set_no
	c.minLen = 0
	c.maxFlag = set_no
	c.maxLen = 0
	c.defaultFlag = set_no
	c.defaultStr = ""
	c.requireFlag = set_no
	c.require = false
	c.inFlag = set_no
	c.in = nil
	c.RegExFlag = set_no
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
	// string length check
	if c.minFlag == set_yes && l < c.minLen {
		if fix && value.CanSet() {
			value.SetString(c.defaultStr)
		}
		return errors.New(fmt.Sprintf("`%s` at least length %d, current length is %d, %s", name, c.minLen, l, v))
	}

	if c.maxFlag == set_yes && l > c.maxLen {
		if fix && value.CanSet() {
			value.SetString(c.defaultStr)
		}
		return errors.New(fmt.Sprintf("`%s` at most length %d, current length is %d, %s", name, c.maxLen, l, v))
	}
	// require
	if c.requireFlag == set_yes && l == 0 {
		if fix && value.CanSet() {
			value.SetString(c.defaultStr)
		}
		return errors.New(fmt.Sprintf("`%s` is empty", name))
	} else if c.requireFlag == set_no && l == 0 {
		// empty string
		return nil
	}
	// in options
	if c.inFlag == set_yes && !inSlice(c.in, v) {
		if fix && value.CanSet() {
			value.SetString(c.defaultStr)
		}
		return errors.New(fmt.Sprintf("`%s` value '%s' is not valid, should be in options", name, v))
	}
	// regex
	if c.RegExFlag == set_yes && len(c.RegExCompile.FindString(v)) == 0 {
		if fix && value.CanSet() {
			value.SetString(c.defaultStr)
		}
		return errors.New(fmt.Sprintf("`%s` value '%s' doesn't match '%s'", name, v, c.RegEx))
	}

	return nil
}

func (c *constraintString) isSet() bool {
	return c.minFlag == set_yes || c.maxFlag == set_yes || c.defaultFlag == set_yes ||
		c.requireFlag == set_yes || c.RegExFlag == set_yes || c.inFlag == set_yes
}

func describeString(fi *reflect.StructField) (constraint, []error) {
	c := new(constraintString)
	c.reset()
	es := make([]error, 0, 0)
	//c.k =  fi.Type.Kind()
	c.k = getLastKind(fi.Type)
	if _rangeIntMap[c.k] == nil {
		return nil, []error{errors.New(fmt.Sprintf("`%s` type is %v, required signed number type", fi.Name, c.k))}
	}
	c.fi = fi
	if minV := fi.Tag.Get(flagMin); len(minV) > 0 {
		if v, e := strconv.ParseInt(minV, 10, 32); e != nil {
			es = append(es, e)
		} else {
			c.minLen = int(v)
			c.minFlag = set_yes
		}
	}
	if maxV := fi.Tag.Get(flagMax); len(maxV) > 0 {
		if v, e := strconv.ParseInt(maxV, 10, 32); e != nil {
			es = append(es, e)
		} else {
			c.maxLen = int(v)
			c.maxFlag = set_yes
		}
	}
	if defaultV := fi.Tag.Get(flagDefault); len(defaultV) > 0 {
		c.defaultStr = defaultV
		c.defaultFlag = set_yes
	}
	if in := fi.Tag.Get(flagIn); len(in) > 0 {
		c.inFlag = set_yes
		c.in = stripEmptyString(splitInOptions(in))
	}
	if req := fi.Tag.Get(flagReq); len(req) > 0 {
		c.requireFlag = set_yes
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
		c.RegExFlag = set_yes
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

	if c.minFlag == set_yes && c.maxFlag == set_yes {
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
	if c.RegExFlag == set_yes {
		if r, e := regexp.Compile(c.RegEx); e != nil {
			es = append(es, e)
			c.RegExFlag = set_no
		} else {
			c.RegExCompile = r
		}
	}

	return es
}
