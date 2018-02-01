package checker

import (
	"testing"
)

type TT struct {
	A string `checker:"len 3,reg ^aasddggas$"`
}

func TestA(t *testing.T) {

	d := TT{"aasddggasd"}
	err := CheckAll(d)
	t.Log(err)
}
