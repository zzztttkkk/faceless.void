package i18n

import (
	"context"
	"fmt"

	"github.com/zzztttkkk/faceless.void/internal"
)

type translateData struct {
	lang string
	txt  string
}

type String struct {
	content string
	langs   []translateData
}

var (
	alli18ns = map[string]*String{}
)

func New(txt string) *String {
	ele, ok := alli18ns[txt]
	if ok {
		return ele
	}
	ele = &String{content: txt}
	alli18ns[txt] = ele
	return ele
}

func (str *String) Format(ctx context.Context, args ...any) string {
	lang := ctx.Value(internal.CtxKeyForLanguageKind)
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

func ExportRaw(filename string) {}
