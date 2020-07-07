package go_validator

import (
	"errors"
	"fmt"
	"log"
	"reflect"
)

func ValidObject(obj interface{}, fix bool) error {
	t := reflect.TypeOf(obj)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return errors.New("object must be struct")
	}
	v := reflect.ValueOf(obj).Elem()
	return doValidObject(&v, fix)
}

func doValidObject(v *reflect.Value, fix bool) error {
	if v.Kind() == reflect.Invalid {
		return nil
	}
	// only process struct
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			fi := v.Type().Field(i)
			val := v.Field(i)
			if _, err := checkFieldValue(&fi, &val, fix); err != nil {
				return err
			}
		}
	}
	return nil
}

func checkFieldValue(fi *reflect.StructField, value *reflect.Value, fix bool) (*FiledTagInfo, error) {
	if value.Kind() == reflect.Invalid {
		return nil, nil
	}
	var fvr *FiledTagInfo = nil

	switch value.Kind() {
	case reflect.Ptr:
		v2 := value.Elem()
		if v2.Kind() == reflect.Invalid {
			return nil, nil
		}
		if v2.Kind() == reflect.Struct {
			return nil, doValidObject(&v2, fix)
		}
		return checkFieldValue(fi, &v2, fix)
	case reflect.Struct:
		return nil, doValidObject(value, fix)
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		if fvr = getRangeValue(fi); fvr == nil || !fvr.IsValid() {
			return nil, nil
		}
		n1 := value.Int()
		// 最小值
		if fvr.Min != nil {
			if m1, ok := fvr.Min.(int64); ok && n1 < m1 {
				if !fix {
					return fvr, errors.New(fmt.Sprintf("%s minimum value is %d，current is %d", fi.Name, m1, n1))
				} else {
					defaultV := m1
					if fvr.Default != nil && n1 == 0 {
						if m2, ok := fvr.Default.(int64); ok {
							defaultV = m2
						}
					}
					if value.CanSet() {
						value.SetInt(defaultV)
					} else {
						log.Printf("can not set value to `%s`", fi.Name)
					}
				}
			}
		}
		// 最大值
		if fvr.Max != nil {
			if m2, ok := fvr.Max.(int64); ok && n1 > m2 {
				if !fix {
					return fvr, errors.New(fmt.Sprintf("%s Maximum value is %d，current is %d", fi.Name, m2, n1))
				} else {
					if value.CanSet() {
						value.SetInt(m2)
					} else {
						log.Printf("can not set value to `%s`", fi.Name)
					}
				}
			}
		}
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		if fvr = getRangeUValue(fi); fvr == nil || !fvr.IsValid() {
			return nil, nil
		}
		n1 := value.Uint()
		// 最小值
		if fvr.Min != nil {
			if m1, ok := fvr.Min.(uint64); ok && n1 < m1 {
				if !fix {
					return fvr, errors.New(fmt.Sprintf("%s minimum value is %d，current is %d", fi.Name, m1, n1))
				} else {
					defaultV := m1
					if fvr.Default != nil && n1 == 0 {
						if m2, ok := fvr.Default.(uint64); ok {
							defaultV = m2
						}
					}
					if value.CanSet() {
						value.SetUint(defaultV)
					} else {
						log.Printf("can not set value to `%s`", fi.Name)
					}
				}
			}
		}
		// 最大值
		if fvr.Max != nil {
			if m2, ok := fvr.Max.(uint64); ok && n1 > m2 {
				if !fix {
					return fvr, errors.New(fmt.Sprintf("%s Maximum value is %d，current is %d", fi.Name, m2, n1))
				} else {
					if value.CanSet() {
						value.SetUint(m2)
					} else {
						log.Printf("can not set value to `%s`", fi.Name)
					}
				}
			}
		}
		//case reflect.String:
		//	fvr = getRangeSValue(fi)
	}
	//fmt.Printf("%v: %v\n", fi.Name, *fvr)
	return fvr, nil
}
