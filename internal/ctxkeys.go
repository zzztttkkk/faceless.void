package internal

type ctxKeyType int

const (
	CtxKeyForHttpRequest = ctxKeyType(iota)
	CtxKeyForBindingGetter
	CtxKeyForLanguageKind
)