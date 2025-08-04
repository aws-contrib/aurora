//go:build !goverter

// Package ent provides an extension to the ent package for database operations.
package ent

import (
	"cmp"
	"context"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate go tool goverter gen .
//go:generate go tool counterfeiter -generate
//go:generate go tool gomarkdoc --output README.md .

// GatewayOption represents a gateway option.
type GatewayOption interface {
	// Apply applies the configuration.
	Apply(*pgxpool.Config) error
}

var _ GatewayOption = GatewayOptionFunc(nil)

// GatewayOptionFunc is a function that applies a GatewayOption.
type GatewayOptionFunc func(*pgxpool.Config) error

// Apply applies the GatewayOptionFunc to the Gateway.
func (fn GatewayOptionFunc) Apply(cfg *pgxpool.Config) error {
	return fn(cfg)
}

// WithURL returns the database URL.
func WithURL() string {
	keys := []string{
		os.Getenv("AURORA_API_DATABASE_URL"),
		os.Getenv("DATABASE_URL"),
	}
	return cmp.Or(keys...)
}

// Open opens a database connection to the given URL.
func Open(ctx context.Context, uri string, options ...GatewayOption) (_ Gateway, err error) {
	// parse the connection string
	config, err := pgxpool.ParseConfig(uri)
	if err != nil {
		return nil, err
	}

	// apply the options
	for _, opt := range options {
		if xerr := opt.Apply(config); xerr != nil {
			return nil, xerr
		}
	}

	// create the connection pool
	conn, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}
	// done!
	return New(conn), nil
}

// RunInTx runs the given function in a transaction.
func (x *Queries) RunInTx(ctx context.Context, action QuerierAction) (err error) {
	type Tx interface {
		Begin(context.Context) (pgx.Tx, error)
	}

	var tx pgx.Tx
	// start the transaction
	if conn, ok := x.db.(Tx); ok {
		if tx, err = conn.Begin(ctx); err != nil {
			return err
		}
	}

	if xerr := action.Run(x.WithTx(tx)); xerr != nil {
		_ = tx.Rollback(ctx)
		return xerr
	}

	return tx.Commit(ctx)
}

// Tx returns the underlying transaction interface.
func (x *Queries) Tx() DBTX {
	return x.db
}

// Ping verifies a connection to the database is still alive,
func (x *Queries) Ping(ctx context.Context) error {
	type Pinger interface {
		Ping(context.Context) error
	}

	conn, _ := x.db.(Pinger)
	return conn.Ping(ctx)
}

// Close closes the connection to the database.
func (x *Queries) Close() {
	type Closer interface {
		Close()
	}

	if conn, ok := x.db.(Closer); ok {
		conn.Close()
	}
}
