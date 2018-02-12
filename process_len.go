package checker

import (
	"fmt"
	"reflect"
	"strconv"
)

type ProcessLen struct {
	Origin   string
	Min, Max int64
}

var _ Process = (*ProcessLen)(nil)

// NewProcessLen "len {min} {max}" or "len {equal}"
func NewProcessLen(tags string) (Process, error) {
	tag := regSpace.Split(tags, -1)
	tag = tag[1:]
	switch len(tag) {
	case 1:
		eq, err := strconv.ParseInt(tag[0], 0, 0)
		if err != nil {
			return nil, err
		}
		return &ProcessLen{
			Origin: tags,
			Min:    eq,
		}, nil
	case 2:
		min, err := strconv.ParseInt(tag[0], 0, 0)
		if err != nil {
			return nil, err
		}
		max, err := strconv.ParseInt(tag[1], 0, 0)
		if err != nil {
			return nil, err
		}

		return &ProcessLen{
			Origin: tags,
			Min:    min,
			Max:    max,
		}, nil
	default:
		return nil, fmt.Errorf("failed `%v` : Len parameter number %v", tags, len(tag))
	}
}

func (p *ProcessLen) CheckValue(v reflect.Value) error {
	switch v.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
	default:
		return fmt.Errorf("failed `%v` : len %s", p.Origin, v.Kind().String())
	}

	l := int64(v.Len())
	if p.Min >= p.Max {
		if l != p.Min {
			return fmt.Errorf("failed `%v` : %v != %v", p.Origin, l, p.Min)
		}
	} else {
		if l < p.Min {
			return fmt.Errorf("failed `%v` : %v < %v", p.Origin, l, p.Min)
		} else if l >= p.Max {
			return fmt.Errorf("failed `%v` : %v >= %v", p.Origin, l, p.Max)
		}
	}
	return nil
}
