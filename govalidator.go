package govalidator

func New() *Validator {
	r := new(Validator)
	return r
}

func ValidObject(obj interface{}, fix bool) []error {
	return New().ValidObject(obj, fix)
}

func Validate(obj interface{}) error {
	return New().Validate(obj)
}

func Validates(obj interface{}) []error {
	return New().Validates(obj)
}
