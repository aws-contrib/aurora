// Package resource provides functions to colorize strings for reporting purposes.
package resource

import (
	"text/template"

	"github.com/fatih/color"
)

var FunMap = template.FuncMap{
	"cyan":   color.CyanString,
	"green":  color.HiGreenString,
	"red":    color.HiRedString,
	"yellow": color.YellowString,
}

// Status represents the status of a migration.
type Status struct {
	Current string `json:"Current,omitempty"` // Current migration version
	Next    string `json:"Next,omitempty"`    // Next migration version
	Status  string `json:"Status,omitempty"`  // Status of migration (OK, PENDING)
}
