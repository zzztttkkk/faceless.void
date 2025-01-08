package vld

type _SliceBuidler[T any] struct {
	_CommonBuilder[[]T, _SliceBuidler[T]]
}

func SliceMeta[T any]() *_SliceBuidler[T] {
	obj := &_SliceBuidler[T]{}
	return obj
}

func Slice[T any](ptr *[]T) *_SliceBuidler[T] {
	return SliceMeta[T]().updateptr(ptr)
}

func (builder *_SliceBuidler[T]) Ele(meta *VldFieldMeta) *_SliceBuidler[T] {
	return builder.set("ele", meta)
}

func (builder *_SliceBuidler[T]) MinSize(minl int) *_SliceBuidler[T] {
	return builder.set("mins", minl)
}

func (builder *_SliceBuidler[T]) MaxSize(maxl int) *_SliceBuidler[T] {
	return builder.set("maxs", maxl)
}

func (builder *_SliceBuidler[T]) NoEmpty() *_SliceBuidler[T] {
	return builder.MinSize(1)
}

func (builder *_SliceBuidler[T]) Optional() *_SliceBuidler[T] {
	return builder.optional()
}
