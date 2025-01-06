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

func (builder *_MapBuidler[K, V]) Key(meta *VldFieldMeta) *_MapBuidler[K, V] {
	return builder.set("key", meta)
}

func (builder *_MapBuidler[K, V]) Ele(meta *VldFieldMeta) *_MapBuidler[K, V] {
	return builder.set("ele", meta)
}
