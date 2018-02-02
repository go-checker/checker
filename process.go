package checker

import (
	"reflect"
)

type MakeProcessFunc func(tag string) (Process, error)

type Process interface {
	CheckValue(reflect.Value) error
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

func (p Processs) CheckValueAll(v reflect.Value) error {
	ee := errs{}
	for _, pf := range p {
		err := pf.CheckValue(v)
		if err != nil {
			ee = append(ee, err)
		}
	}
	if len(ee) == 0 {
		return nil
	}
	return ee
}

type ProcessStruct struct {
	Index   int     // 位置
	Process Process // 处理过程
}

var _ Process = (*ProcessStruct)(nil)

func (p *ProcessStruct) CheckValue(v reflect.Value) error {
	return p.Process.CheckValue(v.Field(p.Index))
}
