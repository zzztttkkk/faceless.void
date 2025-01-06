package vld

import (
	"regexp"
	"sync"
)

type stringBuilder struct {
	_CommonBuilder[string, stringBuilder]
}

func StringMeta() *stringBuilder {
	return (&stringBuilder{})
}

func String(ptr *string) *stringBuilder {
	return StringMeta().updateptr(ptr)
}

func (builder *stringBuilder) Regexp(re *regexp.Regexp) *stringBuilder {
	return builder.set("regexp", re)
}

func (builder *stringBuilder) RegexpString(re string) *stringBuilder {
	return builder.Regexp(regexp.MustCompile(re))
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
	return builder.set("minl", minl)
}

func (builder *stringBuilder) MaxLength(maxl int) *stringBuilder {
	return builder.set("maxl", maxl)
}

func (builder *stringBuilder) NoEmpty() *stringBuilder {
	return builder.MinLength(1)
}

func (builder *stringBuilder) Enum(names ...string) *stringBuilder {
	return builder.set("stringranges", names)
}
