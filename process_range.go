package checker

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type ProcessRange struct {
	Origin   string
	Min, Max float64
}

var _ Process = (*ProcessRange)(nil)

// ProcessRange
// "range {min} {max}"
func NewProcessRange(tags string) (Process, error) {
	tag := strings.Split(tags, " ")
	tag = tag[1:]
	switch len(tag) {
	case 2:
		min, err := strconv.ParseFloat(tag[0], 0)
		if err != nil {
			return nil, err
		}

		max, err := strconv.ParseFloat(tag[1], 0)
		if err != nil {
			return nil, err
		}
		return &ProcessRange{
			Origin: tags,
			Min:    min,
			Max:    max,
		}, nil
	default:
		return nil, fmt.Errorf("Checker failed `%v` : Range parameter number %v", tags, len(tag))
	}
}

func (p *ProcessRange) CheckValue(v reflect.Value) error {
	l := 0.0
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		l = float64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		l = float64(v.Uint())
	case reflect.Float32, reflect.Float64:
		l = v.Float()
	default:
		return fmt.Errorf("Checker failed `%v` : range %s", p.Origin, v.Kind().String())
	}

	if p.Min >= p.Max {
		if l != p.Min {
			return fmt.Errorf("Checker failed `%v` :  %v != %v", p.Origin, l, p.Min)
		}
	} else {
		if l < p.Min {
			return fmt.Errorf("Checker failed `%v` : %v < %v", p.Origin, l, p.Min)
		} else if l >= p.Max {
			return fmt.Errorf("Checker failed `%v` : %v >= %v", p.Origin, l, p.Max)
		}
	}
	return nil
}
