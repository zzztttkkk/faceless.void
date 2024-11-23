package fv

import (
	"fmt"
	"testing"
)

type Code int

const (
	_CodeMin = Code(iota - 1)

	CodeA
	CodeB
	CodeC

	_CodeMax
)

func (code Code) String() string {
	switch code {
	case CodeA:
		{
			return "A"
		}
	case CodeB:
		{
			return "B"
		}
	case CodeC:
		{
			return "C"
		}
	default:
		{
			panic("")
		}
	}
}

func TestEnum(t *testing.T) {
	fmt.Println(EnumPairs(_CodeMin, _CodeMax))
}
