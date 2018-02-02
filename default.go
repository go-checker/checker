package checker

var defaults = NewCheckerClassic()

func CheckValue(i interface{}) error {
	return defaults.Check(i)
}
