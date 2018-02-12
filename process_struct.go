package checker

import "reflect"

type processStruct struct {
	Index   int     // 位置
	Process Process // 处理过程
}

var _ Process = (*processStruct)(nil)

func (p *processStruct) CheckValue(v reflect.Value) error {
	err := p.Process.CheckValue(v.Field(p.Index))
	if err != nil {
		return &errorProcessStruct{
			v:     v,
			index: p.Index,
			err:   err,
		}
	}
	return nil
}
