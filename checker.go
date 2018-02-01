package checker

import (
	"errors"
	"reflect"
	"strings"
)

type Checker struct {
	m   map[string]func(v reflect.Value, tag []string) error
	tag string
}

func NewChecker() *Checker {
	return &Checker{
		tag: "checker",
		m:   map[string]func(v reflect.Value, tag []string) error{},
	}
}

func (c *Checker) SetTag(tag string) {
	c.tag = tag
}

func (c *Checker) AddCheck(name string, fun func(v reflect.Value, tag []string) error) error {
	_, ok := c.m[name]
	if ok {
		return errors.New("error: Defined method " + name)
	}
	c.m[name] = fun
	return nil
}

// 检查遇到第一个就结束
func (c *Checker) Check(i interface{}) error {
	return c.check(reflect.ValueOf(i), func(err error) bool {
		return false
	})
}

// 检查出所有的
func (c *Checker) CheckAll(i interface{}) error {
	ee := errs{}
	err := c.check(reflect.ValueOf(i), func(err error) bool {
		ee = append(ee, err)
		return true
	})
	if err != nil {
		ee = append(ee, err)
	}
	switch len(ee) {
	case 0:
		return nil
	case 1:
		return ee[0]
	default:
		return ee
	}
}

func (c *Checker) check(v reflect.Value, f func(error) bool) error {
	err := checkStruct(v, c.tag, func(vv reflect.Value, t string) error {
		for _, v0 := range strings.Split(t, ",") {
			args := []string{}
			for _, v1 := range strings.Split(v0, " ") {
				v1 = strings.TrimSpace(v1)
				if len(v1) != 0 {
					args = append(args, v1)
				}
			}
			if len(args) == 0 {
				continue
			}
			mf := c.m[args[0]]
			if mf == nil {
				return errors.New("error: Undefined method " + args[0])
			}

			vv = reflect.Indirect(vv)
			if err := mf(vv, args[1:]); err != nil && !f(err) {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func checkStruct(v reflect.Value, tag string, fun func(reflect.Value, string) error) error {
	switch v.Kind() {
	case reflect.Struct:
		t := v.Type()
		for i, l := 0, t.NumField(); i != l; i++ {
			ff := t.Field(i)
			vv, ok := ff.Tag.Lookup(tag)
			if !ok {
				return nil
			}
			err := fun(v.Field(i), vv)
			if err != nil {
				return err
			}
		}
		return nil
	case reflect.Ptr:
		return checkStruct(v.Elem(), tag, fun)
	case reflect.Slice, reflect.Array:
		for i, l := 0, v.Len(); i != l; i++ {
			err := checkStruct(v.Index(i), tag, fun)
			if err != nil {
				return err
			}
		}
		return nil
	case reflect.Map:
		ks := v.MapKeys()
		for _, k := range ks {
			err := checkStruct(v.MapIndex(k), tag, fun)
			if err != nil {
				return err
			}
		}
		return nil
	default:
		return nil
	}
	return nil
}
