package vld

import "github.com/zzztttkkk/faceless.void/internal"

type mapBuidler[K comparable, V any] struct {
	commonBuilder[map[K]V, mapBuidler[K, V]]
}

func MapMeta[K comparable, V any]() *mapBuidler[K, V] {
	obj := &mapBuidler[K, V]{}
	return obj
}

func Map[K comparable, V any](ptr *map[K]V) *mapBuidler[K, V] {
	obj := &mapBuidler[K, V]{}
	return obj.updateptr(ptr)
}

func (builder *mapBuidler[K, V]) Key(meta *VldFieldMeta) *mapBuidler[K, V] {
	builder.pairs = append(builder.pairs, internal.PairOf("key", meta))
	return builder
}

func (builder *mapBuidler[K, V]) Ele(meta *VldFieldMeta) *mapBuidler[K, V] {
	builder.pairs = append(builder.pairs, internal.PairOf("ele", meta))
	return builder
}
