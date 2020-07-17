package Core

import (
	"testing"
)

func TestMoveDir(t *testing.T) {
	outer := Outer{Name: "wc"}
	outer.TestFunc = func() *Inner {
		return NewInner(outer)
	}
	inner := outer.TestFunc()
	t.Log(inner)
}

type Inner struct {
	Name  string
	Outer Outer
}

type Outer struct {
	TestFunc func() *Inner
	Name     string
}

func NewInner(outer Outer) *Inner {
	return &Inner{
		Name:  outer.Name,
		Outer: outer,
	}
}
