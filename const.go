package govalidator

type flagSet int

const (
	setNo  flagSet = iota
	setYes flagSet = iota
)

const (
	flagMin     = "min"     // minimum number, or minimum length of string
	flagMax     = "max"     // maximum number, or maximum length of string
	flagDefault = "default" // default value is missing or error
	flagIn      = "in"      // in option list
	flagReq     = "req"     // required field
	flagRegEx   = "regex"   // regular expression
)

// string option list separator placeholder
const (
	strSep1 = "#"
	strSep2 = ","
)
