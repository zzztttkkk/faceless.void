package main

import (
	_ "embed"
	"fmt"
	"go/types"
	"strings"
)

var (
	gkinds = map[string]func(*PkgInfo) error{}
)

func GenerateOnePkg(info *PkgInfo, kind string) error {
	if kind == "*" {
		for kind, gfnc := range gkinds {
			if err := gfnc(info.Copy().Filter(kind)); err != nil {
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

type HttpHandleFunc struct {
	FuncName string
	Methods  string
	Path     string
	ArgTypes []string

	ReturnKind      string
	ReturnValueType string
}

type HttpInfo struct {
	DicImportName   string
	DICValueName    string
	PkgName         string
	Funcs           []HttpHandleFunc
	MiddlewareNames []string
	Errors          []string
}

//go:embed http.gotpl
var HttpTplSource []byte

func generateOnePkgForHttpGinHandler(info *PkgInfo) error {
	renderdata := &HttpInfo{}

	for _, vd := range info.ValDecls {
		isMiddleware := false
		for _, tag := range vd.tags {
			ps := strings.Split(tag, " ")
			if len(ps) < 2 {
				continue
			}
			if ps[1] == "middleware" {
				isMiddleware = true
				break
			}
		}
		if isMiddleware {
			renderdata.MiddlewareNames = append(renderdata.MiddlewareNames, vd.names...)
		}
	}

	for _, fd := range info.FuncDecls {
		if fd.decl.Type.Params == nil {
			renderdata.Errors = append(renderdata.Errors, fmt.Sprintf(`HandleFunc: "%s", has empty params`, fd.decl.Name))
			continue
		}

		type nt struct {
			name string
			tpy  string
		}

		var args []nt
		for _, p := range fd.decl.Type.Params.List {
			for _, name := range p.Names {
				args = append(args, nt{name: name.String(), tpy: types.ExprString(p.Type)})
			}
		}

		var rets []nt
		if fd.decl.Type.Results != nil {
			for _, p := range fd.decl.Type.Results.List {
				for _, name := range p.Names {
					rets = append(rets, nt{name: name.String(), tpy: types.ExprString(p.Type)})
				}
			}
		}

		if args[0].tpy != "context.Context" {
			renderdata.Errors = append(renderdata.Errors, fmt.Sprintf(`HandleFunc: "%s", first param is not "context.Context"`, fd.decl.Name))
			continue
		}

		var finfo = &HttpHandleFunc{
			FuncName: fd.decl.Name.Name,
			Methods:  "get",
			Path:     "/",
		}

		for i := 1; i < len(args); i++ {
			av := args[i]
			if av.tpy[0] != '*' {
				renderdata.Errors = append(renderdata.Errors, fmt.Sprintf(`HandleFunc: "%s", other params must be a point of struct`, fd.decl.Name))
				continue
			}
			finfo.ArgTypes = append(finfo.ArgTypes, av.tpy[1:])
		}

		fmt.Println(finfo)
	}

	return nil
}

func init() {
	gkinds["http"] = generateOnePkgForHttpGinHandler
}
