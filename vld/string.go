package vld

import "regexp"

type _StringVldRule struct {
	minlen *uint64
	maxlen *uint64
	regexp *regexp.Regexp
}

type stringVldOptionsKey int

const (
	stringVldOptionsKeyForMinLen = stringVldOptionsKey(iota)
	stringVldOptionsKeyForMaxLen
	stringVldOptionsKeyForRegexp
)

type _StringVldOptionPair struct {
	key stringVldOptionsKey
	val any
}

type StringVldOptions struct {
	pairs []_StringVldOptionPair
}

func (opts *StringVldOptions) MinLen(v uint64) *StringVldOptions {
	opts.pairs = append(opts.pairs, _StringVldOptionPair{stringVldOptionsKeyForMinLen, v})
	return opts
}

func (opts *StringVldOptions) MaxLen(v uint64) *StringVldOptions {
	opts.pairs = append(opts.pairs, _StringVldOptionPair{stringVldOptionsKeyForMaxLen, v})
	return opts
}

func (opts *StringVldOptions) Regexp(v *regexp.Regexp) *StringVldOptions {
	opts.pairs = append(opts.pairs, _StringVldOptionPair{stringVldOptionsKeyForRegexp, v})
	return opts
}

func (opts *StringVldOptions) RegexpString(v string) *StringVldOptions {
	opts.pairs = append(opts.pairs, _StringVldOptionPair{stringVldOptionsKeyForRegexp, regexp.MustCompile(v)})
	return opts
}

func (opts *StringVldOptions) Finish() func(string) error {
	return func(v string) error {
		return nil
	}
}

func String() *StringVldOptions {
	return &StringVldOptions{}
}

func init() {
	String().MinLen(1).MaxLen(10).Finish()
}
