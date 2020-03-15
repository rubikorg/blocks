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

func (ic intchecker) IsPositive() bool {
	return ic.result && ic.target > 0
}

func (ic intchecker) IsNegative() bool {
	return ic.result && ic.target < 0
}

func (ic intchecker) IsZero() bool {
	return ic.result && ic.target == 0
}

func (ic intchecker) IsInBetween(min, max int) bool {
	return ic.result && ic.target >= min && ic.target <= max
}
