package vld

import "context"

type IValidate interface {
	Validate(context.Context) error
}
