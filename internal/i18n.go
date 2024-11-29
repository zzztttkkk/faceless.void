package internal

import (
	"context"
	"fmt"
)

type translateData struct {
	lang string
	txt  string
}

type I18nString struct {
	content string
	langs   []translateData
}

var (
	alli18ns = map[string]*I18nString{}
)

func NewI18nString(txt string) *I18nString {
	ele, ok := alli18ns[txt]
	if ok {
		return ele
	}
	ele = &I18nString{content: txt}
	alli18ns[txt] = ele
	return ele
}

func (str *I18nString) Format(ctx context.Context, args ...any) string {
	lang := ctx.Value(CtxKeyForLanguageKind)
	if lang == nil {
		return fmt.Sprintf(str.content, args...)
	}
	for _, ele := range str.langs {
		if ele.lang == lang {
			return fmt.Sprintf(ele.txt, args...)
		}
	}
	return fmt.Sprintf(str.content, args...)
}
