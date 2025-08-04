package cmd

import (
	"net/url"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

var rootCtx = &hcl.EvalContext{
	Functions: map[string]function.Function{
		"getenv": function.New(&function.Spec{
			Params: []function.Parameter{
				{
					Name: "name",
					Type: cty.String,
				},
			},
			Type: function.StaticReturnType(cty.String),
			Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
				name := args[0].AsString()
				val := os.Getenv(name)
				return cty.StringVal(val), nil
			},
		}),
		"urlescape": function.New(&function.Spec{
			Params: []function.Parameter{
				{
					Name: "value",
					Type: cty.String,
				},
			},
			Type: function.StaticReturnType(cty.String),
			Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
				s := args[0].AsString()
				escaped := url.PathEscape(s)
				return cty.StringVal(escaped), nil
			},
		}),
	},
}
