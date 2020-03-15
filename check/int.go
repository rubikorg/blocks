package check

type intchecker struct {
	target int
	Result bool
}

func Int(target int) intchecker {
	return intchecker{
		target: target,
		Result: true,
	}
}

func (ic intchecker) IsPositive() intchecker {
	ic.Result = ic.Result && ic.target > 0
	return ic
}

func (ic intchecker) IsNegative() intchecker {
	ic.Result = ic.Result && ic.target < 0
	return ic
}

func (ic intchecker) IsZero() intchecker {
	ic.Result = ic.Result && ic.target == 0
	return ic
}

func (ic intchecker) IsInBetween(min, max int) intchecker {
	ic.Result = ic.Result && ic.target >= min && ic.target <= max
	return ic
}
