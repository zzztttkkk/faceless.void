package vld

import (
	"regexp"

	"github.com/zzztttkkk/faceless.void/internal"
)

type stringBuilder struct {
	commonBuilder[string, stringBuilder]
}

func String() *stringBuilder {
	return (&stringBuilder{})
}

func StringWithPtr(ptr *string) *stringBuilder {
	return String().updateptr(ptr)
}

func (builder *stringBuilder) Regexp(re *regexp.Regexp) *stringBuilder {
	builder.pairs = append(builder.pairs, internal.PairOf("regexp", re))
	return builder
}

func (builder *stringBuilder) RegexpString(re string) *stringBuilder {
	builder.pairs = append(builder.pairs, internal.PairOf("regexp", regexp.MustCompile(re)))
	return builder
}

func (builder *stringBuilder) MinLength(minl int) *stringBuilder {
	builder.pairs = append(builder.pairs, internal.PairOf("minl", minl))
	return builder
}

func (builder *stringBuilder) MaxLength(maxl int) *stringBuilder {
	builder.pairs = append(builder.pairs, internal.PairOf("maxl", maxl))
	return builder
}

func (builder *stringBuilder) NoEmpty() *stringBuilder {
	builder.pairs = append(builder.pairs, internal.PairOf("minl", 1))
	return builder
}

func (builder *stringBuilder) Enum(names ...string) *stringBuilder {
	builder.pairs = append(builder.pairs, internal.PairOf("stringranges", names))
	return builder
}
