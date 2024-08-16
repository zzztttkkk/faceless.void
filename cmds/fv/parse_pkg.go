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

type ValDecl struct {
	names []string
	tags  []string
}

type PkgInfo struct {
	Dir       string
	Name      string
	FuncDecls []*FuncDecl
	ValDecls  []*ValDecl
}

func (info *PkgInfo) Copy() *PkgInfo {
	ins := &PkgInfo{
		Dir:       info.Dir,
		Name:      info.Name,
		FuncDecls: make([]*FuncDecl, len(info.FuncDecls)),
		ValDecls:  make([]*ValDecl, len(info.ValDecls)),
	}
	copy(ins.FuncDecls, info.FuncDecls)
	copy(ins.ValDecls, info.ValDecls)
	return ins
}

func (info *PkgInfo) Filter(prefix string) *PkgInfo {
	fncs := make([]*FuncDecl, 0)
	for _, v := range info.FuncDecls {
		ntags := []string{}
		for _, tag := range v.tags {
			if strings.HasPrefix(tag, prefix) {
				ntags = append(ntags, tag)
			}
		}
		if len(ntags) > 0 {
			v.tags = ntags
			fncs = append(fncs, v)
		}
	}
	info.FuncDecls = fncs

	vals := make([]*ValDecl, 0)
	for _, v := range info.ValDecls {
		ntags := []string{}
		for _, tag := range v.tags {
			if strings.HasPrefix(tag, prefix) {
				ntags = append(ntags, tag)
			}
		}
		if len(ntags) > 0 {
			v.tags = ntags
			vals = append(vals, v)
		}
	}
	info.ValDecls = vals
	return info
}

func ParseOnePkg(dir string) (*PkgInfo, error) {
	des, e := os.ReadDir(dir)
	if e != nil {
		return nil, e
	}

	pinfo := &PkgInfo{
		Dir: dir,
	}

	fset := token.NewFileSet()
	for _, d := range des {
		if !d.Type().IsRegular() {
			continue
		}
		if !strings.HasSuffix(d.Name(), ".go") || strings.HasSuffix(d.Name(), "_test.go") {
			continue
		}
		if strings.HasPrefix(d.Name(), "fv.generated.") {
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
			switch node := decl.(type) {
			case *ast.FuncDecl:
				{
					fd := node
					if fd.Doc == nil || len(fd.Doc.List) < 1 {
						continue
					}
					fdc := &FuncDecl{
						decl: fd,
					}

					for _, comment := range fd.Doc.List {
						txt := comment.Text
						if strings.HasPrefix(txt, "//go:fv ") {
							fdc.tags = append(fdc.tags, strings.TrimSpace(txt[8:]))
						} else if strings.HasPrefix(txt, "// go:fv ") {
							fdc.tags = append(fdc.tags, strings.TrimSpace(txt[9:]))
						}
					}

					if len(fdc.tags) < 1 {
						continue
					}
					pinfo.FuncDecls = append(pinfo.FuncDecls, fdc)
				}
			case *ast.GenDecl:
				{
					for _, spec := range node.Specs {
						switch spec := spec.(type) {
						case *ast.ValueSpec:
							{
								if spec.Comment == nil {
									continue
								}

								vd := &ValDecl{}
								for _, ident := range spec.Names {
									vd.names = append(vd.names, ident.String())
								}

								for _, c := range spec.Comment.List {
									txt := c.Text
									if strings.HasPrefix(txt, "//go:fv ") {
										vd.tags = append(vd.tags, strings.TrimSpace(txt[8:]))
									} else if strings.HasPrefix(txt, "// go:fv ") {
										vd.tags = append(vd.tags, strings.TrimSpace(txt[9:]))
									}
								}

								if len(vd.tags) < 1 {
									continue
								}

								pinfo.ValDecls = append(pinfo.ValDecls, vd)
							}
						}
					}
				}
			}
		}
	}
	return pinfo, nil
}
