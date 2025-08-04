// Package cmd provides functionality to parse and evaluate HCL configuration files.
package cmd

import (
	"encoding"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"
)

type Node interface {
	// Eval evaluates the node in the context of the provided evaluation context.
	Eval(ctx *hcl.EvalContext) (cty.Value, hcl.Diagnostics)
}

// Config is the top-level struct to hold all parsed HCL blocks.
// It uses struct tags to map HCL blocks to Go fields.
type Config struct {
	Data         []*Data        `hcl:"data,block"`
	Variables    []*Variable    `hcl:"variable,block"`
	Environments []*Environment `hcl:"env,block"`
}

// GetEnvironment retrieves an environment by its name from the Config.
func (c *Config) GetEnvironment(name string) *Environment {
	for _, e := range c.Environments {
		if e.Name == name {
			return e
		}
	}

	return nil
}

var _ encoding.TextUnmarshaler = (*Config)(nil)

// UnmarshalText implements encoding.TextUnmarshaler.
func (c *Config) UnmarshalText(text []byte) error {
	parser := hclparse.NewParser()

	file, ferr := parser.ParseHCL(text, "aurora")
	if ferr.HasErrors() {
		return ferr
	}

	var config Config
	if derr := gohcl.DecodeBody(file.Body, nil, &config); derr.HasErrors() {
		return derr
	}

	config.Eval(rootCtx)

	*c = config
	return nil
}

var _ Node = &Config{}

// Eval implements Node.
func (c *Config) Eval(ctx *hcl.EvalContext) (cty.Value, hcl.Diagnostics) {
	ctx = &hcl.EvalContext{
		Functions: ctx.Functions,
		Variables: map[string]cty.Value{
			"env":  cty.EmptyObjectVal,
			"var":  cty.EmptyObjectVal,
			"data": cty.EmptyObjectVal,
		},
	}

	for _, v := range c.Variables {
		kv := GetValueMap(ctx.Variables["var"])
		kv[v.Name], _ = v.Eval(ctx)
		// Update the context with the evaluated variable
		ctx.Variables["var"] = cty.ObjectVal(kv)
	}

	for _, d := range c.Data {
		kv := GetValueMap(ctx.Variables["data"])
		kv[d.Type], _ = d.Eval(ctx)
		// Update the context with the evaluated variable
		ctx.Variables["data"] = cty.ObjectVal(kv)
	}

	for _, e := range c.Environments {
		kv := GetValueMap(ctx.Variables["env"])
		kv[e.Name], _ = e.Eval(ctx)
		// Update the context with the evaluated variable
		ctx.Variables["env"] = cty.ObjectVal(kv)
	}

	return cty.ObjectVal(ctx.Variables), nil
}

// Environment represents an 'env' block.
type Environment struct {
	Name      string           `hcl:"name,label"`
	Migration *Migration       `hcl:"migration,block"`
	URL       hcl.Expression   `hcl:"url"`
	Context   *hcl.EvalContext `hcl:"-"`
}

// GetURL returns the URL of the environment.
func (x *Environment) GetURL() (string, error) {
	value, err := x.URL.Value(x.Context)
	if err != nil {
		return "", err
	}
	return value.AsString(), nil
}

var _ Node = &Environment{}

// Eval implements Node.
func (x *Environment) Eval(ctx *hcl.EvalContext) (cty.Value, hcl.Diagnostics) {
	x.Context = ctx

	url, err := x.URL.Value(ctx)
	if err != nil {
		return url, err
	}

	migration, err := x.Migration.Eval(ctx)
	if err != nil {
		return cty.Value{}, err
	}

	return cty.ObjectVal(map[string]cty.Value{
		"migration": migration,
		"url":       url,
	}), nil
}

// Migration represents a 'migration' block.
type Migration struct {
	Dir     hcl.Expression   `hcl:"dir"`
	Context *hcl.EvalContext `hcl:"-"`
}

// GetDir returns the directory for migrations.
func (x *Migration) GetDir() (string, error) {
	value, err := x.Dir.Value(x.Context)
	if err != nil {
		return "", err
	}
	return value.AsString(), nil
}

var _ Node = &Migration{}

// Eval implements Node.
func (x *Migration) Eval(ctx *hcl.EvalContext) (cty.Value, hcl.Diagnostics) {
	x.Context = ctx

	dir, err := x.Dir.Value(ctx)
	if err != nil {
		return dir, err
	}

	return cty.ObjectVal(map[string]cty.Value{
		"dir": dir,
	}), nil
}

// Data represents a 'data' block.
type Data struct {
	Type    string           `hcl:"type,label"`
	Name    string           `hcl:"name,label"`
	Remain  hcl.Body         `hcl:",remain"`
	Context *hcl.EvalContext `hcl:"-"`
}

var _ Node = &Data{}

// Eval evaluates the variable's default expression in the context of the provided evaluation context.
func (x *Data) Eval(ctx *hcl.EvalContext) (cty.Value, hcl.Diagnostics) {
	x.Context = ctx
	kv := GetValueMap(ctx.Variables["data"])
	attr := GetValueMap(kv[x.Type])

	switch x.Type {
	case "aws_dsql_token":
		token := &DSQLToken{}
		// decode the token
		if err := gohcl.DecodeBody(x.Remain, ctx, token); err != nil {
			return cty.Value{}, err
		}
		// eval the token
		value, err := token.Eval(ctx)
		if err != nil {
			return cty.Value{}, err
		}

		// done
		attr[x.Name] = value
	default:
		return cty.Value{}, hcl.Diagnostics{
			{
				Severity:    hcl.DiagError,
				Summary:     "The data type is not supported",
				EvalContext: ctx,
			},
		}
	}
	// return the new value
	return cty.ObjectVal(attr), nil
}

// Variable represents a 'variable' block.
type Variable struct {
	Name    string           `hcl:"name,label"`
	Type    hcl.Expression   `hcl:"type"`
	Default hcl.Expression   `hcl:"default,optional"`
	Context *hcl.EvalContext `hcl:"-"`
}

var _ Node = &Variable{}

// Eval evaluates the variable's default expression in the context of the provided evaluation context.
func (x *Variable) Eval(ctx *hcl.EvalContext) (cty.Value, hcl.Diagnostics) {
	x.Context = ctx
	return x.Default.Value(ctx)
}

// GetValueMap retrieves the value map from a cty.Value.
func GetValueMap(v cty.Value) map[string]cty.Value {
	if v.IsNull() {
		v = cty.EmptyObjectVal
	}

	kv := v.AsValueMap()
	if kv == nil {
		kv = make(map[string]cty.Value)
	}

	return kv
}

// GetPath parses a URL and returns the path without the scheme.
func GetPath(value string) (string, error) {
	path, err := url.Parse(value)
	if err != nil {
		return "", err
	}

	switch path.Scheme {
	case "file":
		// remove the scheme
		path.Scheme = ""
	default:
		return "", fmt.Errorf("unsupported config scheme: %s", path.Scheme)
	}

	return strings.TrimPrefix(path.String(), "//"), nil
}
