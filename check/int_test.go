package check

import "testing"

func TestInt(t *testing.T) {
	a := Int(2)
	if a.target != 2 {
		t.Error("Target not set properly when calling Int() method")
	}
}

func TestIntchecker_IsNegative(t *testing.T) {
	if Int(2).IsNegative().Result == true {
		t.Error("Given target is 2 and it is deciding that 2 is a negative number")
	}

	if Int(-2).IsNegative().Result == false {
		t.Error("Given target is -2 and it is deciding that -2 is a positive number")
	}
}

func TestIntchecker_IsPositive(t *testing.T) {
	if Int(892).IsPositive().Result == false {
		t.Error("Given target is 892 and it is deciding that 892 is a negative number")
	}

	if Int(-28829).IsPositive().Result == true {
		t.Error("Given target is -28829 and it is deciding that -28829 is a positive number")
	}
}

func TestIntchecker_IsZero(t *testing.T) {
	if Int(0).IsZero().Result == false {
		t.Error("Given target is 0 and it is deciding that 0 is not 0")
	}

	if Int(-120).IsZero().Result == true {
		t.Error("Given target is -120 and it is deciding that -120 is 0")
	}
}

func TestIntchecker_IsInBetween(t *testing.T) {
	if Int(6).IsInBetween(2, 7).Result == false {
		t.Error("Given target is 6 and it is deciding that 6 is not in between 2 and 7")
	}

	if Int(-86).IsInBetween(-100, -86).Result == false {
		t.Error("Given target is -86 and it is deciding that -86 is not in between -100 and -86")
	}

	if Int(29).IsInBetween(2, 9).Result == true {
		t.Error("Given target is 29 and it is deciding that 29 is in between 2 and 9")
	}
}
