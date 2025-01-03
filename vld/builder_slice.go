package vld

import "github.com/zzztttkkk/faceless.void/internal"

type sliceBuidler[T any] struct {
	commonBuilder[[]T, sliceBuidler[T]]
}

func SliceMeta[T any]() *sliceBuidler[T] {
	obj := &sliceBuidler[T]{}
	return obj
}

func Slice[T any](ptr *[]T) *sliceBuidler[T] {
	return SliceMeta[T]().updateptr(ptr)
}

func (builder *sliceBuidler[T]) Ele(meta *VldFieldMeta) *sliceBuidler[T] {
	builder.pairs = append(builder.pairs, internal.PairOf("ele", meta))
	return builder
}

func (builder *sliceBuidler[T]) MinSize(minl int) *sliceBuidler[T] {
	builder.pairs = append(builder.pairs, internal.PairOf("minl", minl))
	return builder
}

func (builder *sliceBuidler[T]) MaxSize(maxl int) *sliceBuidler[T] {
	builder.pairs = append(builder.pairs, internal.PairOf("maxl", maxl))
	return builder
}

func (builder *sliceBuidler[T]) NoEmpty() *sliceBuidler[T] {
	builder.pairs = append(builder.pairs, internal.PairOf("minl", 1))
	return builder
}
