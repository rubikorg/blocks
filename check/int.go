package check

type intchecker struct {
	target int
	result bool
}

func Int(target int) intchecker {
	return intchecker{
		target: target,
		result: true,
	}
}

func (ic intchecker) IsPositive() intchecker {
	ic.result = ic.result && ic.target > 0
	return ic
}

func (ic intchecker) IsNegative() intchecker {
	ic.result = ic.result && ic.target < 0
	return ic
}

func (ic intchecker) IsZero() intchecker {
	ic.result = ic.result && ic.target == 0
	return ic
}

func (ic intchecker) IsInBetween(min, max int) intchecker {
	ic.result = ic.result && ic.target >= min && ic.target <= max
	return ic
}
