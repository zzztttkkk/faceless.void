package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type FuncDecl struct {
	decl *ast.FuncDecl
	tags []string
}

type PkgInfo struct {
	Name      string
	FuncDecls []*FuncDecl
}

func (info *PkgInfo) Copy() *PkgInfo {
	ins := &PkgInfo{
		Name:      info.Name,
		FuncDecls: make([]*FuncDecl, len(info.FuncDecls)),
	}
	copy(ins.FuncDecls, info.FuncDecls)
	return ins
}

func ParseOnePkg(dir string) (*PkgInfo, error) {
	des, e := os.ReadDir(dir)
	if e != nil {
		return nil, e
	}

	pinfo := &PkgInfo{}

	fset := token.NewFileSet()
	for _, d := range des {
		if !d.Type().IsRegular() {
			continue
		}
		if !strings.HasSuffix(d.Name(), ".go") || strings.HasSuffix(d.Name(), "_test.go") {
			continue
		}

		fp := filepath.Join(dir, d.Name())
		fast, err := parser.ParseFile(fset, fp, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}

		if pinfo.Name == "" {
			pinfo.Name = fast.Name.String()
		} else {
			if pinfo.Name != fast.Name.String() {
				return nil, fmt.Errorf(`package name diff: "%s" "%s"`, pinfo.Name, fast.Name.String())
			}
		}

		for _, decl := range fast.Decls {
			fd, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			if fd.Doc == nil || len(fd.Doc.List) < 1 {
				continue
			}
			fdc := &FuncDecl{
				decl: fd,
			}

			for _, comment := range fd.Doc.List {
				txt := comment.Text
				if strings.HasPrefix(txt, "// :fv:") {
					fdc.tags = append(fdc.tags, txt[7:])
				}
			}

			if len(fdc.tags) < 1 {
				continue
			}
			pinfo.FuncDecls = append(pinfo.FuncDecls, fdc)
		}
	}
	return pinfo, nil
}
