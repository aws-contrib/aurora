//go:build !goverter

package ent

import (
	"context"
	"fmt"
	"io/fs"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

var asyncre = regexp.MustCompile(`(?i)CONCURRENTLY`)

//counterfeiter:generate -o ./fake . FileSystem

// Mutex is a special revision that represents a lock on the database.
var Mutex = &Revision{
	ID:          "20060102150405",
	Description: "lock",
}

// FileSystem represents a filesystem that supports globbing and reading files.
type FileSystem interface {
	fs.FS
	fs.GlobFS
	fs.ReadFileFS
}

// RevisionRepository represents a repository for managing revisions.
type RevisionRepository struct {
	// Gateway represents the database gateway.
	Gateway Gateway
	// FileSystem is the filesystem where the revision files are located.
	FileSystem fs.FS
}

// LockRevisionParams represents the parameters for locking a revision.
type LockRevisionParams struct {
	// Revision contains the parameters for locking a revision.
	Revision *Revision
	// Timeout is the maximum time to wait for the lock.
	Timeout time.Duration
}

// LockRevision locks the revision for exclusive access.
func (x *RevisionRepository) LockRevision(ctx context.Context, params *LockRevisionParams) error {
	start := time.Now()

	for {
		// prepare the parameters
		params.Revision.ExecutedAt = time.Now().UTC()
		params.Revision.ExecutionTime = time.Since(start)

		// create the revision
		args := &InsertRevisionParams{}
		args.SetRevision(params.Revision)
		_, err := x.Gateway.InsertRevision(ctx, args)

		switch {
		case err == nil:
			return nil
		case IsErrorCode(err, ErrCodeUniqueViolation):
			time.Sleep(1 * time.Second)
		case time.Since(start) > params.Timeout:
			return fmt.Errorf("timeout while waiting for lock")
		default:
			return err
		}
	}
}

// UnlockRevisionParams represents the parameters for unlocking a revision.
type UnlockRevisionParams struct {
	// Revision contains the parameters for unlocking a revision.
	Revision *Revision
}

// UnlockRevision unlocks the revision after exclusive access.
func (x *RevisionRepository) UnlockRevision(ctx context.Context, params *UnlockRevisionParams) error {
	args := &ExecDeleteRevisionParams{}
	args.ID = params.Revision.ID
	return x.Gateway.ExecDeleteRevision(ctx, args)
}

// ApplyRevisionParams represents the parameters for executing a revision.
type ApplyRevisionParams struct {
	// Revision contains the parameters for executing a revision.
	Revision *Revision
}

// ApplyRevision executes a revision.
func (x *RevisionRepository) ApplyRevision(ctx context.Context, params *ApplyRevisionParams) error {
	// read the revision content
	data, err := fs.ReadFile(x.FileSystem, params.Revision.GetName())
	if err != nil {
		return err
	}

	queries := strings.SplitAfter(string(data), ";")

	fmt.Println(queries, len(queries))
	// prepare the revision
	params.Revision.Total = len(queries)

	args := &UpsertRevisionParams{}
	args.SetRevision(params.Revision)
	// prepare the revision
	revision, err := x.Gateway.UpsertRevision(ctx, args)
	if err != nil {
		return err
	}

	params.Revision = revision
	// We should skip the already executed statements
	queries = queries[params.Revision.Count:]

	start := time.Now()
	// Apply the statements one by one
	for index, query := range queries {
		query = asyncre.ReplaceAllString(query, "ASYNC")

		job := &Job{}
		// execute the revision
		err = x.Gateway.Tx().QueryRow(ctx, query).Scan(&job.ID)

		switch {
		case err == pgx.ErrNoRows:
			// We are good to go because the operation does not return job id
		case err != nil:
			msg := err.Error()
			// set the revision error
			params.Revision.Error = &msg
		default:
			repository := &JobRepository{
				Gateway: x.Gateway,
			}

			args := &WaitJobParams{}
			args.ID = job.ID
			// Wait for the job to complete
			job, err = repository.WaitJob(ctx, args)
			switch {
			case err == pgx.ErrNoRows:
			case err != nil:
				msg := err.Error()
				// set the revision error
				params.Revision.Error = &msg
			case job.Status == "failed":
				// set the revision error
				params.Revision.Error = &job.Details
			}
		}

		params.Revision.Count = index + 1
		params.Revision.ErrorStmt = &query
		params.Revision.ExecutedAt = time.Now().UTC()
		params.Revision.ExecutionTime = time.Since(start)

		args := &ExecUpdateRevisionParams{}
		args.SetRevision(params.Revision)

		// prepare the mask
		args.UpdateMask = append(args.UpdateMask, "executed_at")
		args.UpdateMask = append(args.UpdateMask, "execution_time")

		if args.Error == nil {
			args.UpdateMask = append(args.UpdateMask, "count")
		} else {
			args.UpdateMask = append(args.UpdateMask, "error")
			args.UpdateMask = append(args.UpdateMask, "error_stmt")
		}

		if err := x.Gateway.ExecUpdateRevision(ctx, args); err != nil {
			return err
		}
	}

	return nil
}

// ListRevisions lists all revisions in the repository.
func (x *RevisionRepository) ListRevisions(ctx context.Context, _ *ListRevisionsParams) (collection []*Revision, _ error) {
	matches, err := fs.Glob(x.FileSystem, "*.sql")
	if err != nil {
		return nil, err
	}

	for index, path := range sort.StringSlice(matches) {
		revision := &Revision{}
		revision.SetName(path)

		// read the revision content
		data, err := fs.ReadFile(x.FileSystem, revision.GetName())
		if err != nil {
			return nil, err
		}

		queries := strings.SplitAfter(string(data), ";")
		// prepare the revision
		revision.Total = len(queries)

		// append the revision
		collection = append(collection, revision)

		params := &GetRevisionParams{}
		params.SetRevision(revision)

		// load the revision
		revision, err = x.Gateway.GetRevision(ctx, params)
		switch {
		case err == pgx.ErrNoRows:
		// We are good to go because the revision does not exist
		case err != nil:
			return nil, err
		default:
			collection[index] = revision
		}
	}

	return collection, nil
}
