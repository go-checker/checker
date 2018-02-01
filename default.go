package checker

var defaults = NewCheckerClassic()

func Check(i interface{}) error {
	return defaults.Check(i)
}

func CheckAll(i interface{}) error {
	return defaults.CheckAll(i)
}

func NewCheckerClassic() *Checker {
	c := NewChecker()
	c.AddCheck("len", FuncLen)
	c.AddCheck("regexp", FuncRegexp)
	c.AddCheck("range", FuncRange)
	return c
}
