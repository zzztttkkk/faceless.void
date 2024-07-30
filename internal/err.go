package internal

import (
	"fmt"
	"os"
)

type _ErrNamespace struct {
	ns string
}

func (e _ErrNamespace) Errorf(tpl string, args ...any) error {
	return fmt.Errorf(fmt.Sprintf(`fv.%s: %s`, e.ns, tpl), args...)
}

func ErrNamespace(ns string) _ErrNamespace {
	return _ErrNamespace{ns: ns}
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
