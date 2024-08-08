package main

import (
	_ "embed"
	"fmt"
)

var (
	gkinds = map[string]func(*PkgInfo) error{}
)

func GenerateOnePkg(info *PkgInfo, kind string) error {
	if kind == "*" {
		for _, gfnc := range gkinds {
			if err := gfnc(info.Copy()); err != nil {
				return err
			}
		}
		return nil
	}
	gfnc, ok := gkinds[kind]
	if !ok {
		return fmt.Errorf("unknown generate kind: %s", kind)
	}
	return gfnc(info)
}

type HttpInfo struct {
	DicImportName string
	DICValueName  string
	PkgName       string
	Funcs         []struct {
		FuncName string
		Methods  string
		Path     string
		Args     []struct {
			StructType string
		}
		ReturnKind      string
		ReturnValueType string
	}
}

//go:embed http.gotpl
var HttpTplSource []byte

func generateOnePkgForHttpGinHandler(info *PkgInfo) error {
	fmt.Println(info.Name)
	return nil
}

func init() {
	gkinds["http"] = generateOnePkgForHttpGinHandler
}
