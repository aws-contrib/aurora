package ent_test

import (
	"context"

	"github.com/aws-contrib/aurora/internal/database/ent"
)

// NewGateway returns a new test gateway.
func NewGateway() (ent.Gateway, error) {
	return ent.Open(context.TODO(), ent.WithURL())
}
