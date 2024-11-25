package vld

import (
	"fmt"
	"regexp"
	"slices"

	"github.com/zzztttkkk/faceless.void/internal"
)

type stringVldOptionsKey int

const (
	stringVldOptionsKeyForMinLen = stringVldOptionsKey(iota)
	stringVldOptionsKeyForMaxLen
	stringVldOptionsKeyForRegexp
	stringVldOptionsKeyForEnum
	stringVldOptionsKeyForCustom
)

type _VldOptionPair[K ~int] struct {
	key K
	val any
}

type _StringVldOptionPair = _VldOptionPair[stringVldOptionsKey]

type _StringVldOptions struct {
	pairs []_StringVldOptionPair
}

func (opts *_StringVldOptions) MinLen(v int) *_StringVldOptions {
	opts.pairs = append(opts.pairs, _StringVldOptionPair{stringVldOptionsKeyForMinLen, v})
	return opts
}

func (opts *_StringVldOptions) MaxLen(v int) *_StringVldOptions {
	opts.pairs = append(opts.pairs, _StringVldOptionPair{stringVldOptionsKeyForMaxLen, v})
	return opts
}

func (opts *_StringVldOptions) Regexp(v *regexp.Regexp) *_StringVldOptions {
	opts.pairs = append(opts.pairs, _StringVldOptionPair{stringVldOptionsKeyForRegexp, v})
	return opts
}

func (opts *_StringVldOptions) RegexpString(v string) *_StringVldOptions {
	opts.pairs = append(opts.pairs, _StringVldOptionPair{stringVldOptionsKeyForRegexp, regexp.MustCompile(v)})
	return opts
}

func (opts *_StringVldOptions) Enum(names []string) *_StringVldOptions {
	opts.pairs = append(opts.pairs, _StringVldOptionPair{stringVldOptionsKeyForEnum, names})
	return opts
}

func (opts *_StringVldOptions) Custom(fnc func(string) error) *_StringVldOptions {
	opts.pairs = append(opts.pairs, _StringVldOptionPair{stringVldOptionsKeyForCustom, fnc})
	return opts
}

func (opts *_StringVldOptions) Func() func(string) error {
	var fncs []func(string) error

	for _, pair := range opts.pairs {
		switch pair.key {
		case stringVldOptionsKeyForMinLen:
			{
				minlen := pair.val.(int)
				if minlen >= 0 {
					fncs = append(fncs, func(s string) error {
						if len(s) < minlen {
							return fmt.Errorf("")
						}
						return nil
					})
				}
				break
			}
		case stringVldOptionsKeyForMaxLen:
			{
				maxlen := pair.val.(int)
				if maxlen >= 0 {
					fncs = append(fncs, func(s string) error {
						if len(s) > maxlen {
							return fmt.Errorf("")
						}
						return nil
					})
				}
				break
			}
		case stringVldOptionsKeyForRegexp:
			{
				regexp := pair.val.(*regexp.Regexp)
				if regexp != nil {
					fncs = append(fncs, func(s string) error {
						if !regexp.MatchString(s) {
							return fmt.Errorf("")
						}
						return nil
					})
				}
				break
			}
		case stringVldOptionsKeyForEnum:
			{
				names := pair.val.([]string)
				if len(names) < 16 {
					fncs = append(fncs, func(s string) error {
						if !slices.Contains(names, s) {
							return fmt.Errorf("")
						}
						return nil
					})
				} else {
					set := make(internal.Set[string])
					for _, txt := range names {
						set[txt] = internal.Empty{}
					}

					fncs = append(fncs, func(s string) error {
						_, ok := set[s]
						if !ok {
							return fmt.Errorf("")
						}
						return nil
					})
				}
			}
		case stringVldOptionsKeyForCustom:
			{
				fnc := pair.val.(func(string) error)
				if fnc != nil {
					fncs = append(fncs, fnc)
				}
				break
			}
		}
	}

	if len(fncs) < 1 {
		return nil
	}
	return func(v string) error {
		for _, fnc := range fncs {
			err := fnc(v)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func String() *_StringVldOptions {
	return &_StringVldOptions{}
}

type _Strings struct{}

var (
	Strings _Strings
)

// https://github.com/go-playground/validator/blob/master/regexes.go
func (*_Strings) Email() *_StringVldOptions {
	return String().RegexpString("^(?:(?:(?:(?:[a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22))))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$")
}

func (*_Strings) ObjectId() *_StringVldOptions {
	return String().RegexpString("^[a-f\\d]{24}$")
}

func (*_Strings) Hex() *_StringVldOptions {
	return String().RegexpString("^(0[xX])?[0-9a-fA-F]+$")
}