package vld

import (
	"context"
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

type _StringVldOptions struct {
	pairs []internal.Pair[stringVldOptionsKey]
}

func (opts *_StringVldOptions) MinLen(v int) *_StringVldOptions {
	opts.pairs = append(opts.pairs, internal.PairOf(stringVldOptionsKeyForMinLen, v))
	return opts
}

func (opts *_StringVldOptions) MaxLen(v int) *_StringVldOptions {
	opts.pairs = append(opts.pairs, internal.PairOf(stringVldOptionsKeyForMaxLen, v))
	return opts
}

func (opts *_StringVldOptions) Regexp(v *regexp.Regexp) *_StringVldOptions {
	opts.pairs = append(opts.pairs, internal.PairOf(stringVldOptionsKeyForRegexp, v))
	return opts
}

func (opts *_StringVldOptions) RegexpString(v string) *_StringVldOptions {
	opts.pairs = append(opts.pairs, internal.PairOf(stringVldOptionsKeyForRegexp, regexp.MustCompile(v)))
	return opts
}

func (opts *_StringVldOptions) Enum(names []string) *_StringVldOptions {
	opts.pairs = append(opts.pairs, internal.PairOf(stringVldOptionsKeyForEnum, names))
	return opts
}

func (opts *_StringVldOptions) Custom(fnc func(string) error) *_StringVldOptions {
	opts.pairs = append(opts.pairs, internal.PairOf(stringVldOptionsKeyForCustom, fnc))
	return opts
}

var (
	msgForStringLengthLtMin    = internal.NewI18nString(`fv.vld: string length less than min(%d)`)
	msgForStringLengthGtMax    = internal.NewI18nString(`fv.vld: string length greater than max(%d)`)
	msgForStringNotMatchRegexp = internal.NewI18nString(`fv.vld: string not match regexp`)
	msgForStringNotInEnums     = internal.NewI18nString(`fv.vld: string not in enums`)
)

func (opts *_StringVldOptions) Func() func(context.Context, string) error {
	var fncs []func(context.Context, string) error

	for _, pair := range opts.pairs {
		switch pair.Key {
		case stringVldOptionsKeyForMinLen:
			{
				minlen := pair.Val.(int)
				if minlen >= 0 {
					fncs = append(fncs, func(ctx context.Context, s string) error {
						if len(s) < minlen {
							return internal.NewError(ctx, ErrorKindVldStringLengthLtMin, msgForStringLengthLtMin, minlen)
						}
						return nil
					})
				}
				break
			}
		case stringVldOptionsKeyForMaxLen:
			{
				maxlen := pair.Val.(int)
				if maxlen >= 0 {
					fncs = append(fncs, func(ctx context.Context, s string) error {
						if len(s) > maxlen {
							return internal.NewError(ctx, ErrorKindVldStringLengthGtMax, msgForStringLengthGtMax, maxlen)
						}
						return nil
					})
				}
				break
			}
		case stringVldOptionsKeyForRegexp:
			{
				regexp := pair.Val.(*regexp.Regexp)
				if regexp != nil {
					fncs = append(fncs, func(ctx context.Context, s string) error {
						if !regexp.MatchString(s) {
							return internal.NewError(ctx, ErrorKindVldStringNotMatchRegexp, msgForStringNotMatchRegexp)
						}
						return nil
					})
				}
				break
			}
		case stringVldOptionsKeyForEnum:
			{
				names := pair.Val.([]string)
				if len(names) < 16 {
					fncs = append(fncs, func(ctx context.Context, s string) error {
						if !slices.Contains(names, s) {
							return internal.NewError(ctx, ErrorKindVldStringNotInEnums, msgForStringNotInEnums)
						}
						return nil
					})
				} else {
					set := make(internal.Set[string])
					for _, txt := range names {
						set[txt] = internal.Empty{}
					}

					fncs = append(fncs, func(ctx context.Context, s string) error {
						_, ok := set[s]
						if !ok {
							return internal.NewError(ctx, ErrorKindVldStringNotInEnums, msgForStringNotInEnums)
						}
						return nil
					})
				}
			}
		case stringVldOptionsKeyForCustom:
			{
				fnc := pair.Val.(func(context.Context, string) error)
				if fnc != nil {
					fncs = append(fncs, func(ctx context.Context, s string) error {
						err := fnc(ctx, s)
						if err != nil {
							return internal.NewError(ctx, ErrorKindVldCustomFunc, msgForCustomFunc, err, s)
						}
						return nil
					})
				}
				break
			}
		}
	}

	if len(fncs) < 1 {
		return nil
	}
	return func(ctx context.Context, v string) error {
		for _, fnc := range fncs {
			err := fnc(ctx, v)
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
