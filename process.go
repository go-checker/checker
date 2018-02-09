package checker

import (
	"reflect"
)

type MakeProcessFunc func(tag string) (Process, error)

type Process interface {
	CheckValue(reflect.Value) error
}

type processMaps map[reflect.Value]Processs

func (p processMaps) check() error {
	for k, pf := range p {
		err := pf.CheckValue(k)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p processMaps) checkAll() error {
	ee := errs{}
	for k, pf := range p {
		for _, v := range pf {
			err := v.CheckValue(k)
			if err != nil {
				ee = append(ee, err)
			}
		}
	}
	if len(ee) == 0 {
		return nil
	}
	return ee
}

type Processs []Process

var _ Process = (*Processs)(nil)

func (p Processs) CheckValue(v reflect.Value) error {
	for _, pf := range p {
		err := pf.CheckValue(v)
		if err != nil {
			return err
		}
	}
	return nil
}
