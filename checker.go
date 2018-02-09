package checker

import (
	"errors"
	"reflect"
	"strings"
)

type Checker struct {
	mf  map[string]MakeProcessFunc // Make process function
	mp  map[string]Processs        // Cache processs
	tag string                     // tag
	all bool                       // check all
}

var _ Process = (*Checker)(nil)

func NewChecker() *Checker {
	c := &Checker{
		tag: "checker",
		mf:  map[string]MakeProcessFunc{},
		mp:  map[string]Processs{},
	}
	c.AddCheck("check", c.processCheck)
	c.AddCheck("len", NewProcessLen)
	c.AddCheck("range", NewProcessRange)
	c.AddCheck("regexp", NewProcessRegexp)
	return c
}

func NewCheckerAll() *Checker {
	c := NewChecker()
	c.all = true
	return c
}

func (c *Checker) SetTag(tag string) {
	c.tag = tag
}

func (c *Checker) AddCheck(name string, fun MakeProcessFunc) error {
	_, ok := c.mf[name]
	if ok {
		return errors.New("error: Defined method " + name)
	}
	c.mf[name] = fun
	return nil
}

func (c *Checker) Check(i interface{}) error {
	return c.CheckValue(reflect.ValueOf(i))
}

func (c *Checker) CheckValue(v reflect.Value) error {
	ms := processMaps{}
	err := c.process(v, ms)
	if err != nil {
		return err
	}

	if c.all {
		return ms.checkAll()
	}
	return ms.check()
}

func (c *Checker) parserTag(tag string) (prs Processs, err error) {
	for _, v := range strings.Split(tag, ",") {
		v = strings.TrimSpace(v)
		i := strings.Index(v, " ")
		k := v
		if i > 0 {
			k = k[:i]
		}
		mp, ok := c.mf[k]
		if !ok {
			return nil, errors.New("error: Undefined method " + k)
		}
		pr, err := mp(v)
		if err != nil {
			return nil, err
		}
		prs = append(prs, pr)
	}

	return prs, nil
}

func (c *Checker) processStruct(t reflect.Type) (prs Processs, err error) {
	ppn := t.PkgPath() + "." + t.Name()
	if tp, ok := c.mp[ppn]; ok {
		return tp, nil
	}

	for i, l := 0, t.NumField(); i != l; i++ {
		tf := t.Field(i)

		// 获取成员的tag
		tv, ok := tf.Tag.Lookup(c.tag)
		if !ok {
			continue
		}

		if len(tv) == 0 {
			continue
		}

		if tv[0] == '-' {
			continue
		}

		// 拼凑 唯一key
		mk := ppn + "." + tf.Name

		// 尝试获取已经解析过的
		pp, ok := c.mp[mk]
		if !ok {
			// 生成tag 解析
			pp, err = c.parserTag(tv)
			if err != nil {
				return nil, err
			}
			c.mp[mk] = pp
		}

		// 记录结果
		for _, p := range pp {
			prs = append(prs, &ProcessStruct{
				Index:   i,
				Process: p,
			})
		}
	}
	c.mp[ppn] = prs
	return prs, nil
}

// 找出所有要处理的 value 和 process
func (c *Checker) process(v reflect.Value, ms processMaps) (err error) {
	switch v.Kind() {
	case reflect.Struct:
		t := v.Type()
		prs, err := c.processStruct(t)
		if err != nil {
			return err
		}
		if len(prs) == 0 {
			return nil
		}
		ms[v] = prs
	case reflect.Ptr:
		return c.process(v.Elem(), ms)
	case reflect.Slice, reflect.Array:
		for i, l := 0, v.Len(); i != l; i++ {
			err := c.process(v.Index(i), ms)
			if err != nil {
				return err
			}
		}
		return nil
	case reflect.Map:
		for _, k := range v.MapKeys() {
			err := c.process(v.MapIndex(k), ms)
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
