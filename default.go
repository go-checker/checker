package checker

var defaults = NewChecker()

func Check(i interface{}) error {
	return defaults.Check(i)
}
