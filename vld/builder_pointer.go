package vld

type _PointerBuilder[T any] struct {
	_CommonBuilder[T, _PointerBuilder[T]]
}

func PointerMeta[T any]() *_PointerBuilder[*T] {
	obj := &_PointerBuilder[*T]{}
	return obj
}

func Pointer[T any](ptr **T) *_PointerBuilder[*T] {
	return PointerMeta[T]().updateptr(ptr)
}

func (builder *_PointerBuilder[T]) Optional() *_PointerBuilder[T] {
	return builder.optional()
}

func (builder *_PointerBuilder[T]) Ele(meta *VldFieldMeta) *_PointerBuilder[T] {
	return builder.set("ele", meta)
}
