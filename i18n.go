package fv

import "github.com/zzztttkkk/faceless.void/internal"

type I18nString internal.I18nString

func NewI18nString(v string) *I18nString {
	return (*I18nString)(internal.NewI18nString(v))
}
