package checker

import (
	"reflect"
)

type processCheck struct {
	Checker *Checker
}

var _ Process = (*processCheck)(nil)

// processCheck
// "check"
func (c *Checker) processCheck(tags string) (Process, error) {
	return &processCheck{
		Checker: c,
	}, nil
}

func (p *processCheck) CheckValue(v reflect.Value) error {
	return p.Checker.CheckValue(v)
}
