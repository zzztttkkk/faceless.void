package vld

type _MapBuidler[K comparable, V any] struct {
	_CommonBuilder[map[K]V, _MapBuidler[K, V]]
}

func MapMeta[K comparable, V any]() *_MapBuidler[K, V] {
	obj := &_MapBuidler[K, V]{}
	return obj
}

func Map[K comparable, V any](ptr *map[K]V) *_MapBuidler[K, V] {
	obj := &_MapBuidler[K, V]{}
	return obj.updateptr(ptr)
}

func (builder *_MapBuidler[K, V]) Optional() *_MapBuidler[K, V] {
	return builder.optional()
}

func (builder *_MapBuidler[K, V]) MinSize(minl int) *_MapBuidler[K, V] {
	return builder.set("mins", minl)
}

func (builder *_MapBuidler[K, V]) MaxSize(maxl int) *_MapBuidler[K, V] {
	return builder.set("maxs", maxl)
}

func (builder *_MapBuidler[K, V]) NoEmpty() *_MapBuidler[K, V] {
	return builder.MinSize(1)
}

func (builder *_MapBuidler[K, V]) Key(meta *VldFieldMeta) *_MapBuidler[K, V] {
	return builder.set("key", meta)
}

func (builder *_MapBuidler[K, V]) Ele(meta *VldFieldMeta) *_MapBuidler[K, V] {
	return builder.set("ele", meta)
}
