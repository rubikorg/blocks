package checker

import (
	"fmt"
	"strconv"
)

func IsZero(val interface{}) error {
	err := IsStr(val)
	if err != nil {
		return err
	}

	tgt, err := strconv.Atoi(val.(string))
	if err != nil {
		return fmt.Errorf("%v cannot be converted into an integer", val)
	}

	if tgt != 0 {
		return fmt.Errorf("value %v does not equate to zero", val)
	}

	return nil
}
