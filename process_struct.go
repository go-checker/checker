package checker

import "reflect"

type ProcessStruct struct {
	Index   int     // 位置
	Process Process // 处理过程
}

var _ Process = (*ProcessStruct)(nil)

func (p *ProcessStruct) CheckValue(v reflect.Value) error {
	return p.Process.CheckValue(v.Field(p.Index))
}
