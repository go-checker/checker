package checker

import (
	"fmt"
	"reflect"
	"strings"
)

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

type errorProcessStruct struct {
	v     reflect.Value
	index int
	err   error
}

func (e *errorProcessStruct) Error() string {
	t := e.v.Type()
	fi := t.Field(e.index)
	return fmt.Sprintf("%s.%s: %s", t.Name(), fi.Name, e.err)
}
