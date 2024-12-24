package internalvld

import (
	"github.com/zzztttkkk/lion"
)

func Ptr[T any]() *T {
	return lion.Ptr[T]()
}

func FieldOf[T any](ptr any) *lion.Field[VldFieldMeta] {
	return lion.FieldOf[T, VldFieldMeta](ptr)
}
