package internal

type ctxKeyType int

const (
	CtxKeyForAppScope = ctxKeyType(iota)
	CtxKeyForVldScheme
	CtxKeyForLanguageKind
)
