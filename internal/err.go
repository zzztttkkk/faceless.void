package internal

import "fmt"

type _ErrNamespace struct {
	ns string
}

func (e _ErrNamespace) Errorf(tpl string, args ...any) error {
	return fmt.Errorf(fmt.Sprintf(`fv.%s: %s`, e.ns, tpl), args...)
}

func ErrNamespace(ns string) _ErrNamespace {
	return _ErrNamespace{ns: ns}
}
