package checker

import "strings"

type errs []error

func (e errs) String() string {
	s := []string{"", "Checker failed:"}
	for _, v := range e {
		s = append(s, v.Error())
	}
	s = append(s, "")
	return strings.Join(s, "\n")
}

func (e errs) Error() string {
	return e.String()
}
