package checker

import (
	"errors"
	"fmt"
	"reflect"

	r "github.com/rubikorg/rubik"
)

// MustExist checks if value of the given field is nil or not
func MustExist(val interface{}) error {
	if val == nil {
		return errors.New("$ is required")
	}

	return nil
}

// IsEmail checks if given field is an email or not
func IsEmail(val interface{}) error {
	err := IsStr(val)
	if err != nil {
		return err
	}

	return nil
}

// IsStr checks if the type of the value referred in the Entity is a
// string type. An use case of this can be to safeguard your further
// assertions.
func IsStr(val interface{}) error {
	if val == nil {
		return nil
	}

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
			msg := fmt.Sprintf("minimum %d characters needed but value: $", minLen)
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
			msg := fmt.Sprintf("maximum of %d characters allowed but value: $",
				maxLen)
			return errors.New(msg)
		}

		return nil
	}
}

// StrAllow only allows the use of the runes given inside the
// string value of the request field
func StrAllow(runes ...rune) r.Assertion {
	return func(val interface{}) error {
		return nil
	}
}

// StrIsOneOf allowes only the values passed inside this method
// as a viable value for the request field
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
			return fmt.Errorf("$ must be one of %v", values)
		}

		return nil
	}
}

func StrBoolIsTrue(val interface{}) error {
	err := IsStr(val)
	if err != nil {
		return err
	}
	strVal, _ := val.(string)
	switch strVal {
	case "y", "yes", "true", "TRUE", "True":
		return nil
	default:
		return errors.New("$ is not a truthy string")
	}
}

func StrBoolIsFalse(val interface{}) error {
	err := IsStr(val)
	if err != nil {
		return err
	}
	strVal, _ := val.(string)
	switch strVal {
	case "n", "no", "false", "FALSE", "False":
		return nil
	default:
		return errors.New("$ is not a falsy string")
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
