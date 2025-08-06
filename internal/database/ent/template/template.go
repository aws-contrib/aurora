// Package template provides functions to colorize strings for reporting purposes.
package template

import (
	"embed"
	"fmt"
	"io"
	"text/template"

	"github.com/fatih/color"
)

var FunMap = template.FuncMap{
	"cyan":   color.CyanString,
	"green":  color.HiGreenString,
	"red":    color.HiRedString,
	"yellow": color.YellowString,
	"sub": func(x, y int) int {
		return x - y
	},
}

//go:embed *.txt
var fs embed.FS

// Execute executes a template with the given name and data, writing the output to the provided writer.
func Execute(w io.Writer, name string, data any) error {
	name = fmt.Sprintf("template_%s.txt", name)

	tmpl, err := template.New(name).Funcs(FunMap).ParseFS(fs, name)
	if err != nil {
		return err
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		return err
	}
	return nil
}
