package checker

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type ProcessRegexp struct {
	Origin string
	Reg    *regexp.Regexp
}

var _ Process = (*ProcessRegexp)(nil)

// NewProcessRegexp "regexp {text}"
func NewProcessRegexp(tags string) (Process, error) {
	i := strings.Index(tags, " ")
	expr := tags[i+1:]
	reg, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	return &ProcessRegexp{
		Origin: tags,
		Reg:    reg,
	}, nil
}

func (p *ProcessRegexp) CheckValue(v reflect.Value) error {
	s := v.String()
	if !p.Reg.MatchString(s) {
		return fmt.Errorf("invalid `%v`: %s", p.Origin, s)
	}
	return nil
}
