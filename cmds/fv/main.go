package main

import (
	"github.com/alecthomas/kong"
)

type App struct {
	Generate struct {
		Recurse   bool     `help:"recurse scan"`
		ScanPaths []string `arg:"" name:"paths" help:"paths to scan."`
	} `cmd:""`
}

func main() {
	app := &App{}
	ctx := kong.Parse(app)
	switch ctx.Command() {
	case "generate <paths>":
		{
			RunGenerate(app.Generate.ScanPaths)
			return
		}
	default:
		{
			panic(ctx.Command())
		}
	}
}
