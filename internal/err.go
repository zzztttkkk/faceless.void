package internal

import (
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
