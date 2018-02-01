package checker

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// FuncRegexp
// "regexp {text}"
func FuncRegexp(v reflect.Value, tag []string) (err error) {
	pat := strings.Join(tag, " ")
	s := v.String()
	matched, err := regexp.MatchString(pat, s)
	if err != nil {
		return err
	}
	if !matched {
		return fmt.Errorf("checker reg `%v`: %s", pat, s)
	}
	return nil
}

// FuncRange
// "range {min} {max}"
func FuncRange(v reflect.Value, tag []string) (err error) {
	if len(tag) != 2 {
		return fmt.Errorf("checker range %v: checker: Len parameter number %v", tag, len(tag))
	}
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		l := v.Int()
		min, err := strconv.ParseInt(tag[0], 0, 0)
		if err != nil {
			return err
		}
		if l < min {
			return fmt.Errorf("checker range %v: %v < %v", tag, l, tag[0])
		}
		max, err := strconv.ParseInt(tag[1], 0, 0)
		if err != nil {
			return err
		}
		if l >= max {
			return fmt.Errorf("checker range %v: %v >= %v", tag, l, tag[1])
		}
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		l := v.Uint()
		min, err := strconv.ParseUint(tag[0], 0, 0)
		if err != nil {
			return err
		}
		if l < min {
			return fmt.Errorf("checker range %v: %v < %v", tag, l, tag[0])
		}
		max, err := strconv.ParseUint(tag[1], 0, 0)
		if err != nil {
			return err
		}
		if l >= max {
			return fmt.Errorf("checker range %v: %v >= %v", tag, l, tag[1])
		}
		return nil
	case reflect.Float32, reflect.Float64:
		l := v.Float()
		min, err := strconv.ParseFloat(tag[0], 0)
		if err != nil {
			return err
		}
		if l < min {
			return fmt.Errorf("checker range %v: %v < %v", tag, l, tag[0])
		}
		max, err := strconv.ParseFloat(tag[1], 0)
		if err != nil {
			return err
		}
		if l >= max {
			return fmt.Errorf("checker range %v: %v >= %v", tag, l, tag[1])
		}
		return nil
	default:
		return fmt.Errorf("checker range %v: range %s", tag, v.Kind().String())
	}
	return nil
}

// LenFunc
// "len {min} {max}" or "len {equal}"
func FuncLen(v reflect.Value, tag []string) (err error) {

	l := int64(-1)
	switch v.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		l = int64(v.Len())
	default:
		return fmt.Errorf("checker len %v: len %s", tag, v.Kind().String())
	}

	switch len(tag) {
	case 1:
		eq, err := strconv.ParseInt(tag[0], 0, 0)
		if err != nil {
			return err
		}
		if l != eq {
			return fmt.Errorf("checker len %v:  %v != %v", tag, l, tag[0])
		}
		return nil
	case 2:
		min, err := strconv.ParseInt(tag[0], 0, 0)
		if err != nil {
			return err
		}
		if l < min {
			return fmt.Errorf("checker len %v: %v < %v", tag, l, tag[0])
		}
		max, err := strconv.ParseInt(tag[1], 0, 0)
		if err != nil {
			return err
		}
		if l >= max {
			return fmt.Errorf("checker len %v: %v >= %v", tag, l, tag[1])
		}
		return nil
	default:
		return fmt.Errorf("checker len %v: checker: Len parameter number %v", tag, len(tag))
	}
	return nil
}
