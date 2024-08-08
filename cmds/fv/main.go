package main

import (
	"fmt"
	"github.com/alecthomas/kong"
)

type App struct {
	Generate struct {
		Recurse      bool     `help:"recurse scan"`
		ScanPatterns []string `arg:"" name:"patterns" help:"path patterns to scan."`
	} `cmd:""`
}

func main() {
	app := &App{}
	ctx := kong.Parse(app)
	switch ctx.Command() {
	case "generate":
		{
			fmt.Println("Generating scan patterns...", app.Generate)
		}
	default:
		{
			panic(ctx.Command())
		}
	}
}
