package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func collectpks(p string, visited map[string]struct{}, pkgs *[]string) {
	p, e := filepath.Abs(p)
	if e != nil {
		return
	}

	stat, e := os.Stat(p)
	if e != nil {
		return
	}
	if !stat.IsDir() {
		return
	}
	if stat.Mode()&fs.ModeSymlink != 0 {
		return
	}

	_, ok := visited[p]
	if ok {
		return
	}
	visited[p] = struct{}{}

	des, e := os.ReadDir(p)
	if e != nil {
		return
	}
	ispkg := false
	for _, d := range des {
		if strings.HasSuffix(d.Name(), ".go") && d.Type().IsRegular() {
			ispkg = true
			break
		}
	}
	if ispkg {
		*pkgs = append(*pkgs, p)
	}

	for _, d := range des {
		if d.IsDir() {
			collectpks(fmt.Sprintf("%s/%s", p, d.Name()), visited, pkgs)
		}
	}
}

func RunGenerate(paths []string) {
	var pkgs []string
	visited := map[string]struct{}{}
	for _, path := range paths {
		collectpks(path, visited, &pkgs)
	}

	for _, pkg := range pkgs {
		e := DoGenerateOnePkg(pkg)
		if e != nil {
			panic(e)
		}
	}
}

func DoGenerateOnePkg(dir string) error {
	info, e := ParseOnePkg(dir)
	if e != nil {
		return e
	}
	return GenerateOnePkg(info, "*")
}
