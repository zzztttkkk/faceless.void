package vld

import (
	"regexp"
	"sync"

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

var (
	emialRegexpString = ("^(?:(?:(?:(?:[a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22))))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$")
	getEmialRegexp    = func() func() *regexp.Regexp {
		var regex *regexp.Regexp
		var once sync.Once
		return func() *regexp.Regexp {
			once.Do(func() {
				regex = regexp.MustCompile(emialRegexpString)
			})
			return regex
		}
	}()
)

func (builder *stringBuilder) Email() *stringBuilder {
	return builder.Regexp(getEmialRegexp())
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
