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

type comparator func(int64, int64) bool

func checkMinMax(f reflect.StructField, v reflect.Value, m int, c comparator) bool {
	switch f.Type.Kind() {
	case reflect.Int:
		if !c(int64(m), v.Int()) {
			return false
		}
	case reflect.String:
		if !c(int64(m), int64(len(v.String()))) {
			return false
		}
	}
	return true
}

func contains(f reflect.StructField, v reflect.Value, c []string, e *ValidationErrors) bool {
	for _, s := range c {
		switch f.Type.Kind() {
		case reflect.Int:
			n, err := strconv.Atoi(s)
			if err != nil {
				*e = append(*e, ValidationError{ErrInvalidValidatorSyntax})
				return false
			}
			if v.Int() == int64(n) {
				return true
			}
		case reflect.String:
			if v.String() == s {
				return true
			}
		}
	}
	return false
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
				if len(vv.Field(i).String()) != val {
					errs = append(errs, ValidationError{
						fmt.Errorf("field %s has an invalid length", field.Name)})
				}

			case "min":
				c := func(min int64, v int64) bool {
					return v >= min
				}

				if !checkMinMax(field, vv.Field(i), val, c) {
					errs = append(errs, ValidationError{
						fmt.Errorf("field %s has value less than min", field.Name)})
				}

			case "max":
				c := func(max int64, v int64) bool {
					return v <= max
				}

				if !checkMinMax(field, vv.Field(i), val, c) {
					errs = append(errs, ValidationError{
						fmt.Errorf("field %s has value bigger than max", field.Name)})
				}
			case "in":
				if !contains(field, vv.Field(i), strings.Split(toCheck, ","), &errs) {
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
