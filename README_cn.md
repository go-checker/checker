# Checker
golang 检查结构体里的字段数据是否正确

 - [English](./README.md)
 - [简体中文](./README_cn.md)
 
## 安装
``` bash
go get -u -v github.com/wzshiming/checker
```

## 用法

[API 文档](http://godoc.org/github.com/wzshiming/checker)

[示例](./examples/main.go)


``` golang
package main

import (
	"gopkg.in/ffmt.v1"

	checker "github.com/wzshiming/checker"
)

type TT0 struct {
	B string `checker:"len 5"`
	C int    `checker:"range 5 19"`
}

func main() {
	checkAll := checker.NewCheckerAll()
	check := checker.NewChecker()
	var err error
	d0 := []struct {
		A  string `checker:"len 4,regexp d"`
		Ts []TT0  `checker:"len 2"`
	}{
		{
			A: "ssa",
		},
	}
	err = check.Check(d0)
	if err != nil {
		ffmt.Mark(err)
		ffmt.Puts(d0)
	}
	/*
		main.go:28  .A: failed `len 4` : 3 != 4
		[
		 {
		  A:  "ssa"
		  Ts: [ ]
		 }
		]
	*/

	err = checkAll.Check(d0)
	if err != nil {
		ffmt.Mark(err)
		ffmt.Puts(d0)
	}
	/*
		main.go:43
		Checker failed:
		.A: failed `len 4` : 3 != 4
		.A: failed `regexp d` : ssa
		.Ts: failed `len 2` : 0 != 2

		[
		 {
		  A:  "ssa"
		  Ts: [ ]
		 }
		]
	*/

	d1 := struct {
		A  string `checker:"len 4,regexp ^ss"`
		Ts []TT0  `checker:"check"`
	}{
		A: "aa",
		Ts: []TT0{
			{
				B: "13",
			},
		},
	}
	err = checkAll.Check(d1)
	if err != nil {
		ffmt.Mark(err)
		ffmt.Puts(d1)
	}
	/*
		main.go:74
		Checker failed:
		.A: failed `len 4` : 2 != 4
		.A: failed `regexp ^ss` : aa
		.Ts:
		Checker failed:
		TT0.B: failed `len 5` : 2 != 5
		TT0.C: failed `range 5 19` : 0 < 5

		{
		 A:  "aa"
		 Ts: [
		  {
		   B: "13"
		   C: 0
		  }
		 ]
		}
	*/
}
```

## MIT许可证

版权所有©2017-2018 wzshiming <[https://github.com/wzshiming](https://github.com/wzshiming)>。

MIT是[MIT许可证](https://opensource.org/licenses/MIT)许可的开源软件。