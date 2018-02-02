package checker

import (
	"testing"

	ffmt "gopkg.in/ffmt.v1"
)

type TT1 struct {
	B string `checker:"len 5"`
	C int    `checker:"range 5 19"`
}
type TT struct {
	A   string `checker:"len 4,regexp d"`
	TT1 []TT1  `checker:"len 2,check"`
}

func TestA(t *testing.T) {

	d := TT{
		A: "ssa",
		TT1: []TT1{
			{
				B: "13",
			},
		},
	}
	ch := NewCheckerAll()
	err := ch.Check(d)
	ffmt.Mark(err)
	//ffmt.Puts(ch.mp)
}
