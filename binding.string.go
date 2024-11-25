package fv

import (
	"context"
	"fmt"
	"strings"
	"unsafe"
)

type _StringFieldBinding struct {
	ins   *_BindingInstance
	name  string
	alias []string
	where BindingSrcKind
	ptr   *string
	vld   func(string) error

	optional       bool
	trimspace      bool
	defaultvalue   string
	defaultvalueok bool
}

func (ins *_BindingInstance) String(ptr *string) *_StringFieldBinding {
	sfb := &_StringFieldBinding{
		ins: ins,
		ptr: ptr,
	}
	ins.fields = append(ins.fields, sfb.do)
	return sfb
}

func (sfb *_StringFieldBinding) Default(v string) *_StringFieldBinding {
	sfb.defaultvalue = v
	sfb.defaultvalueok = true
	return sfb
}

func (sfb *_StringFieldBinding) Optional() *_StringFieldBinding {
	sfb.optional = true
	return sfb
}

func (sfb *_StringFieldBinding) TrimSpace() *_StringFieldBinding {
	sfb.trimspace = true
	return sfb
}

func (sfb *_StringFieldBinding) Name(name string, alias ...string) *_StringFieldBinding {
	sfb.name = name
	sfb.alias = alias
	return sfb
}

func (sfb *_StringFieldBinding) From(src BindingSrcKind) *_StringFieldBinding {
	sfb.where = src
	return sfb
}

func (sfb *_StringFieldBinding) Validate(vld func(string) error) *_StringFieldBinding {
	sfb.vld = vld
	return sfb
}

func (sfb *_StringFieldBinding) do(ctx context.Context) error {
	if sfb.name == "" {
		sfb.name = sfb.ins.nameof(unsafe.Pointer(sfb.ptr))
	}
	val, ok := BindingGetter(ctx).String(sfb.where, sfb.name, sfb.alias...)
	if !ok {
		if !sfb.defaultvalueok {
			if !sfb.optional {
				return fmt.Errorf("missing required filed: `%s`", sfb.name)
			}
			return nil
		}
		val = sfb.defaultvalue
	}
	if sfb.trimspace {
		val = strings.TrimSpace(val)
	}
	if sfb.vld != nil {
		err := sfb.vld(val)
		if err != nil {
			return err
		}
	}
	*sfb.ptr = val
	return nil
}
