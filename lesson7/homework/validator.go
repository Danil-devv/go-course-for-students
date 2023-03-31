package homework

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"strconv"
	"strings"
)

var ErrNotStruct = errors.New("wrong argument given, should be a struct")
var ErrInvalidValidatorSyntax = errors.New("invalid validator syntax")
var ErrValidateForUnexportedFields = errors.New("validation for unexported field is not allowed")

type ValidationError struct {
	Err error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	res := ""
	for _, err := range v {
		res = res + err.Err.Error()
	}
	return res
}

func createReflectionSlice(vv reflect.Value) []reflect.Value {
	v := make([]reflect.Value, 0)

	switch vv.Type().Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < vv.Len(); i++ {
			v = append(v, vv.Index(i))
		}
	default:
		v = append(v, vv)
	}
	return v
}

func checkLen(vv reflect.Value, ln int) bool {
	v := createReflectionSlice(vv)

	for i := 0; i < len(v); i++ {
		if len(v[i].String()) != ln {
			return false
		}
	}
	return true
}

type comparator func(int64, int64) bool

func checkMinMax(vv reflect.Value, m int, c comparator) bool {
	v := createReflectionSlice(vv)

	for i := 0; i < len(v); i++ {
		switch v[i].Kind() {
		case reflect.Int:
			if !c(int64(m), v[i].Int()) {
				return false
			}
		case reflect.String:
			if !c(int64(m), int64(len(v[i].String()))) {
				return false
			}
		}
	}
	return true
}

func checkContains(vv reflect.Value, c []string, e *ValidationErrors) bool {
	v := createReflectionSlice(vv)

	for i := 0; i < len(v); i++ {
		ok := false
		for _, s := range c {
			switch v[i].Kind() {
			case reflect.Int:
				n, err := strconv.Atoi(s)
				if err != nil {
					*e = append(*e, ValidationError{ErrInvalidValidatorSyntax})
					return false
				}
				if v[i].Int() == int64(n) {
					ok = true
				}
			case reflect.String:
				if v[i].String() == s {
					ok = true
				}
			}
		}
		if !ok {
			return false
		}
	}
	return true
}

func Validate(v any) error {
	vt := reflect.TypeOf(v)
	vv := reflect.ValueOf(v)

	if vt.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	errs := make(ValidationErrors, 0)

	for i, field := range reflect.VisibleFields(vt) {
		if s, ok := field.Tag.Lookup("validate"); !field.IsExported() && ok {
			errs = append(errs, ValidationError{ErrValidateForUnexportedFields})
			continue
		} else if ok {

			validator, toCheck := strings.Split(s, ":")[0], strings.Split(s, ":")[1]

			var val int

			switch validator {
			case "len", "min", "max":
				v, err := strconv.Atoi(toCheck)
				val = v
				if err != nil {
					errs = append(errs, ValidationError{ErrInvalidValidatorSyntax})
					continue
				}
			}

			switch validator {
			case "len":
				if !checkLen(vv.Field(i), val) {
					errs = append(errs, ValidationError{
						fmt.Errorf("field %s has an invalid length", field.Name)})
				}

			case "min":
				c := func(min int64, v int64) bool {
					return v >= min
				}

				if !checkMinMax(vv.Field(i), val, c) {
					errs = append(errs, ValidationError{
						fmt.Errorf("field %s has value less than min", field.Name)})
				}

			case "max":
				c := func(max int64, v int64) bool {
					return v <= max
				}

				if !checkMinMax(vv.Field(i), val, c) {
					errs = append(errs, ValidationError{
						fmt.Errorf("field %s has value bigger than max", field.Name)})
				}
			case "in":
				if !checkContains(vv.Field(i), strings.Split(toCheck, ","), &errs) {
					errs = append(errs, ValidationError{
						fmt.Errorf("field %s does not occur in %v", field.Name,
							strings.Split(toCheck, ","))})
				}
			}
		}
	}

	if len(errs) != 0 {
		return errs
	}

	return nil
}
