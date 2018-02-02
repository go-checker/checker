package checker

import (
	"errors"
	"reflect"
	"strings"
)

type Checker struct {
	mf  map[string]MakeProcessFunc // 历遍结构体生成校验函数
	mp  map[string]Processs        // 直接对结构体校验的函数
	tag string                     // tag
	all bool                       // 是检查全部错误
}

var _ Process = (*Checker)(nil)

func NewChecker() *Checker {
	c := &Checker{
		tag: "checker",
		mf:  map[string]MakeProcessFunc{},
		mp:  map[string]Processs{},
	}
	c.AddCheck("check", c.processCheck)
	return c
}

func NewCheckerClassic() *Checker {
	c := NewChecker()
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

func NewCheckerAllClassic() *Checker {
	c := NewCheckerClassic()
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
	switch v.Kind() {
	case reflect.Struct:
		t := v.Type()
		ps, err := c.ProcessStruct(t)
		if err != nil {
			return err
		}
		if c.all {
			return ps.CheckValueAll(v)
		}

		return ps.CheckValue(v)
	case reflect.Ptr:
		return c.CheckValue(v.Elem())
	case reflect.Slice, reflect.Array:
		for i, l := 0, v.Len(); i != l; i++ {
			err := c.CheckValue(v.Index(i))
			if err != nil {
				return err
			}
		}
		return nil
	case reflect.Map:
		ks := v.MapKeys()
		for _, k := range ks {
			err := c.CheckValue(v.MapIndex(k))
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

func (c *Checker) ParserTag(tag string) (prs Processs, err error) {
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

func (c *Checker) ProcessStruct(t reflect.Type) (prs Processs, err error) {
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

		// 拼凑 唯一key
		mk := ppn + "." + tf.Name

		// 尝试获取已经解析过的
		pp, ok := c.mp[mk]
		if !ok {
			// 生成tag 解析
			pp, err = c.ParserTag(tv)
			if err != nil {
				return nil, err
			}
			c.mp[mk] = pp
		}

		// 记录结果
		if c.all {
			for _, p := range pp {
				prs = append(prs, &ProcessStruct{
					Index:   i,
					Process: p,
				})
			}
			continue
		}

		prs = append(prs, &ProcessStruct{
			Index:   i,
			Process: pp,
		})
	}
	c.mp[ppn] = prs
	return prs, nil
}
