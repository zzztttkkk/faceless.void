package internal

import (
	"context"
	"fmt"
	"os"
)

type ErrNamespace struct {
	Namespace string
}

func (e ErrNamespace) Errorf(tpl string, args ...any) error {
	return fmt.Errorf(fmt.Sprintf(`fv.%s: %s`, e.Namespace, tpl), args...)
}

func FsExists(fp string) (bool, error) {
	_, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

type ErrorKind int

const (
	MaxVldErrorKind     ErrorKind = 1000
	MaxBindingErrorKind ErrorKind = 2000
)

type Error struct {
	Kind    ErrorKind
	Args    []any
	Message string
}

func (e Error) Error() string {
	return e.Message
}

func NewError(ctx context.Context, kind ErrorKind, i18n *I18nString, args ...any) error {
	return Error{
		Kind:    kind,
		Message: i18n.Format(ctx, args...),
		Args:    args,
	}
}
