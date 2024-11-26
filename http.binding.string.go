package fv

import (
	"context"
	"fmt"
	"strings"
	"unsafe"
)

type BindingSliceOptions struct {
	MinLength  int
	MaxLength  int
	VldWithIdx bool
}

type _BindingStringField struct {
	ins   *_BindingInstance
	name  string
	alias []string
	where BindingSrcKind

	ptr      *string
	sliceptr *[]string

	vld func(context.Context, string) error

	optional       bool
	trimspace      bool
	defaultvalue   string
	defaultvalueok bool
}

func (ins *_BindingInstance) String(ptr *string) *_BindingStringField {
	sfb := &_BindingStringField{
		ins: ins,
		ptr: ptr,
	}
	ins.fields = append(ins.fields, sfb.do)
	return sfb
}

func (ins *_BindingInstance) StringSlice(ptr *[]string, opts *BindingSliceOptions) *_BindingStringField {
	sfb := &_BindingStringField{
		ins:      ins,
		sliceptr: ptr,
	}
	ins.fields = append(ins.fields, sfb.do)
	return sfb
}

func (sfb *_BindingStringField) Default(v string) *_BindingStringField {
	if sfb.sliceptr != nil {
		panic("fv.binding: can not set default value for slice")
	}
	sfb.defaultvalue = v
	sfb.defaultvalueok = true
	return sfb
}

func (sfb *_BindingStringField) Optional() *_BindingStringField {
	sfb.optional = true
	return sfb
}

func (sfb *_BindingStringField) TrimSpace() *_BindingStringField {
	sfb.trimspace = true
	return sfb
}

func (sfb *_BindingStringField) Name(name string) *_BindingStringField {
	sfb.name = name
	return sfb
}

func (sfb *_BindingStringField) Alias(alias ...string) *_BindingStringField {
	sfb.alias = alias
	return sfb
}

func (sfb *_BindingStringField) From(src BindingSrcKind) *_BindingStringField {
	sfb.where = src
	return sfb
}

func (sfb *_BindingStringField) Validate(vld func(context.Context, string) error) *_BindingStringField {
	sfb.vld = vld
	return sfb
}

func (sfb *_BindingStringField) do(ctx context.Context) error {
	if sfb.name == "" {
		if sfb.ptr != nil {
			sfb.name = sfb.ins.nameof(unsafe.Pointer(sfb.ptr))
		} else {
			sfb.name = sfb.ins.nameof(unsafe.Pointer(sfb.sliceptr))
		}
	}

	if sfb.ptr != nil {
		val, ok := _BindingGetter(ctx).String(sfb.where, sfb.name, sfb.alias...)
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
			err := sfb.vld(ctx, val)
			if err != nil {
				return err
			}
		}
		*sfb.ptr = val
	} else {
		vals, ok := _BindingGetter(ctx).Strings(sfb.where, sfb.name, sfb.alias...)
		if !ok {
			if !sfb.optional {

			}
			return nil
		}

		if sfb.vld == nil {
			for _, val := range vals {
				if sfb.trimspace {
					val = strings.TrimSpace(val)
				}
				*(sfb.sliceptr) = append(*(sfb.sliceptr), val)
			}
		} else {
			for _, val := range vals {
				if sfb.trimspace {
					val = strings.TrimSpace(val)
				}
				err := sfb.vld(ctx, val)
				if err != nil {
					return err
				}
				*(sfb.sliceptr) = append(*(sfb.sliceptr), val)
			}
		}
	}
	return nil
}
