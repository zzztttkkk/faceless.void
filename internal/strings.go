package internal

import "unsafe"

func StringToReadonlyBytes(v string) []byte {
	return unsafe.Slice(unsafe.StringData(v), len(v))
}

func BytesToString(bs []byte) string {
	return unsafe.String(unsafe.SliceData(bs), len(bs))
}
