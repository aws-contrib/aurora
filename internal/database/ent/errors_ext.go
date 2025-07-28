package ent

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	// ErrTooManyRows occurs when more rows than expected are returned.
	ErrTooManyRows = pgx.ErrTooManyRows
	// ErrNoRows occurs when rows are expected but none are returned.
	ErrNoRows = pgx.ErrNoRows
)

// IsErrorNotFound reports whether the error is a "not found" error.
func IsErrorNotFound(err error) bool {
	return errors.Is(err, ErrNoRows)
}

// PgError represents an error reported by the PostgreSQL server.
type Error = pgconn.PgError

var ErrCodeUniqueViolation = pgerrcode.UniqueViolation

// IsErrorCode reports whether the error is a PostgreSQL error with the given code.
func IsErrorCode(err error, code string) bool {
	var pgerr *Error

	if errors.As(err, &pgerr) {
		return pgerr.Code == code
	}

	return false
}
