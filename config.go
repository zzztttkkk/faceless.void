package fv

import (
	"fmt"
	"os"
	"reflect"

	"dario.cat/mergo"
	"github.com/pelletier/go-toml/v2"
)

func readTomlOne(typ reflect.Type, fp string) (any, error) {
	bs, err := os.ReadFile(fp)
	if err != nil {
		return nil, err
	}
	ptr := reflect.New(typ).Interface()
	err = toml.Unmarshal(bs, ptr)
	if err != nil {
		return nil, err
	}
	return ptr, nil
}

func LoadConfig[T any](typehint T, fps ...string) (*T, error) {
	if len(fps) < 1 {
		return nil, fmt.Errorf("empty config file list")
	}

	typ := reflect.TypeOf(typehint)
	var result *T
	for _, fp := range fps {
		val, err := readTomlOne(typ, fp)
		if err != nil {
			return nil, fmt.Errorf("read config file error: %w", err)
		}
		if result == nil {
			result = val.(*T)
		} else {
			err = mergo.Merge(result, val, mergo.WithOverride)
			if err != nil {
				return nil, err
			}
		}
	}
	return result, nil
}
