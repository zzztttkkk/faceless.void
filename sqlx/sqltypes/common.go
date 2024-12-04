package sqltypes

import (
	"unsafe"

	"github.com/zzztttkkk/faceless.void/internal"
)

type typecommon[T any, D any] struct {
	pairs []internal.Pair[string]
}

type SqlTypeArgs struct {
	Kind string
	Args []any
}

func (c *typecommon[T, D]) sqltype(name string, kind string, args ...any) {
	c.pairs = append(c.pairs, internal.PairOf("name", name))
	c.pairs = append(c.pairs, internal.PairOf("sqltype", SqlTypeArgs{Kind: kind, Args: args}))
}

func (c *typecommon[T, D]) self() *T {
	return (*T)(unsafe.Pointer(c))
}

func (c *typecommon[T, D]) Nullable() *T {
	c.pairs = append(c.pairs, internal.PairOf("nullable", true))
	return c.self()
}

func (c *typecommon[T, D]) Primary() *T {
	c.pairs = append(c.pairs, internal.PairOf("primary", true))
	return c.self()
}

func (c *typecommon[T, D]) Unique() *T {
	c.pairs = append(c.pairs, internal.PairOf("unique", true))
	return c.self()
}

func (c *typecommon[T, D]) Default(v D) *T {
	c.pairs = append(c.pairs, internal.PairOf("default", v))
	return c.self()
}

func (c *typecommon[T, D]) DefaultExpr(expr string) *T {
	c.pairs = append(c.pairs, internal.PairOf("defaultexpr", expr))
	return c.self()
}

func (c *typecommon[T, D]) Check(expr string) *T {
	c.pairs = append(c.pairs, internal.PairOf("check", expr))
	return c.self()
}

func (c *typecommon[T, D]) Comment(comment string) *T {
	c.pairs = append(c.pairs, internal.PairOf("comment", comment))
	return c.self()
}
