package fv_test

import (
	"fmt"
	"testing"

	fv "github.com/zzztttkkk/faceless.void"
)

type DBConfig struct {
	URI string `toml:"uri"`
}

type AppConfig struct {
	Debug bool `toml:"debug"`
	DB    struct {
		Master DBConfig   `toml:"master"`
		Slaves []DBConfig `toml:"slaves"`
	} `toml:"db"`
}

func TestLoadConfig(t *testing.T) {
	ptr, err := fv.LoadConfig(AppConfig{}, "./a.toml", "./a.local.toml", "./a.dev.toml")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(ptr)
}
