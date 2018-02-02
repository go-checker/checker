package checker

var defaults = NewChecker()

func CheckValue(i interface{}) error {
	return defaults.Check(i)
}
