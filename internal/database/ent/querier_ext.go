package ent

import (
	"context"

	"github.com/jackc/pgx/v5"
)

//counterfeiter:generate -o ./fake . DBTX
//counterfeiter:generate -o ./fake . Querier
//counterfeiter:generate -o ./fake . Gateway

// Gateway represents the database gateway.
type Gateway interface {
	// inherit from Querier
	Querier
	// RunInTx runs the given function in a transaction.
	RunInTx(context.Context, QuerierAction) error
	// Ping verifies a connection to the database is still alive.
	Ping(context.Context) error
	// Close closes the connection to the database.
	Close()
}

//counterfeiter:generate -o ./fake github.com/jackc/pgx/v5.Tx
//counterfeiter:generate -o ./fake github.com/jackc/pgx/v5.Row
//counterfeiter:generate -o ./fake github.com/jackc/pgx/v5.Rows
//counterfeiter:generate -o ./fake github.com/jackc/pgx/v5.BatchResults

// Batch represents a batch of results.
type Batch struct {
	Results pgx.BatchResults
	Total   int
}

// QuerierAction represents a query action.
type QuerierAction interface {
	// Run runs the action.
	Run(Querier) error
}

var _ QuerierAction = QuerierFunc(nil)

// QuerierFunc is a function that runs a query.
type QuerierFunc func(Querier) error

// Run runs the query.
func (fn QuerierFunc) Run(querier Querier) error {
	return fn(querier)
}

// NewQueryPipeline returns a new QuerierFunc that runs the given steps.
func NewQueryPipeline(collection ...QuerierFunc) QuerierAction {
	fn := func(querier Querier) error {
		for _, action := range collection {
			if err := action.Run(querier); err != nil {
				return err
			}
		}

		return nil
	}

	return QuerierFunc(fn)
}
