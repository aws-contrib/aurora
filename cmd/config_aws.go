package cmd

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dsql/auth"
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

// DSQLToken represents a token used in the configuration.
type DSQLToken struct {
	Username hcl.Expression `hcl:"username"`
	Endpoint hcl.Expression `hcl:"endpoint"`
	Region   hcl.Expression `hcl:"region"`
}

func (x *DSQLToken) Eval(ctx *hcl.EvalContext) (cty.Value, hcl.Diagnostics) {
	username, err := x.Username.Value(ctx)
	if err != nil {
		return cty.Value{}, err
	}

	endpoint, err := x.Endpoint.Value(ctx)
	if err != nil {
		return cty.Value{}, err
	}

	region, err := x.Region.Value(ctx)
	if err != nil {
		return cty.Value{}, err
	}

	cfg, xerr := config.LoadDefaultConfig(context.Background(), config.WithRegion(region.AsString()))

	if xerr != nil {
		return cty.Value{}, hcl.Diagnostics{
			{
				Severity: hcl.DiagError,
				Summary:  "Unable to load AWS configuration",
				Detail:   xerr.Error(),
			},
		}
	}

	var authz DSQLAuthFunc

	if username.AsString() == "admin" {
		authz = auth.GenerateDBConnectAdminAuthToken
	} else {
		authz = auth.GenerateDbConnectAuthToken
	}

	// build token
	token, terr := authz(context.Background(), endpoint.AsString(), cfg.Region, cfg.Credentials)
	if terr != nil {
		return cty.Value{}, hcl.Diagnostics{
			{
				Severity: hcl.DiagError,
				Summary:  "Unable to generate DSQL authentication token",
				Detail:   terr.Error(),
			},
		}
	}

	return cty.StringVal(token), nil
}

// DSQLAuthFunc is a function that generates an authentication token for AWS DSQL.
type DSQLAuthFunc func(ctx context.Context, endpoint, region string, creds aws.CredentialsProvider, optFns ...func(options *auth.TokenOptions)) (string, error)
