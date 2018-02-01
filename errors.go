package checker

import "strings"

type errs []error

func (e errs) String() string {
	s := []string{}
	for _, v := range e {
		s = append(s, v.Error())
	}
	return strings.Join(s, "\n")
}

func (e errs) Error() string {
	return e.String()
}
