package checker

import (
	"errors"
	"fmt"
	"reflect"

	r "github.com/rubikorg/rubik"
)

func IsEmail(val interface{}) error {
	if _, ok := val.(string); !ok {
		return errors.New("checker(IsEmail): %s cannot be asserted as a string value")
	}
	return nil
}

// IsString checks if the type of the value referred in the Entity is a
// string type. An use case of this can be to safeguard your further
// assertions.
func IsStr(val interface{}) error {
	isReflectedAsString := reflect.TypeOf(val).Kind() == reflect.String
	if isReflectedAsString {
		return nil
	}

	if _, ok := val.(string); !ok {
		return errors.New("cannot be asserted as string")
	}

	return nil
}

func StrMin(minLen int) r.Assertion {
	return func(val interface{}) error {
		err := IsStr(val)
		if err != nil {
			return err
		}

		if len(val.(string)) < minLen {
			msg := fmt.Sprintf("minimum %d characters needed but value: %s", minLen, val.(string))
			return errors.New(msg)
		}

		return nil
	}
}

func StrMax(maxLen int) r.Assertion {
	return func(val interface{}) error {
		err := IsStr(val)
		if err != nil {
			return err
		}

		if len(val.(string)) > maxLen {
			msg := fmt.Sprintf("maximum of %d characters allowed but value: %s",
				maxLen, val.(string))
			return errors.New(msg)
		}

		return nil
	}
}

func StrIsOneOf(values ...string) r.Assertion {
	return func(val interface{}) error {
		if len(values) == 0 {
			return nil
		}

		var vals []interface{}
		for _, a := range values {
			vals = append(vals, a)
		}

		ok := isOneOf(val, vals...)
		if !ok {
			return fmt.Errorf("%v must be one of %v", val, values)
		}

		return nil
	}
}

func isOneOf(t interface{}, values ...interface{}) bool {
	for _, v := range values {
		if t == v {
			return true
		}
	}

	return false
}
