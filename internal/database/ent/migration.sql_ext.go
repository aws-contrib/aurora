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

var (
	commentRegexp      = regexp.MustCompile(`(?m)^\s*--.*$`)
	concurrentlyRegexp = regexp.MustCompile(`(?i)CONCURRENTLY`)
)

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

// MigrationRepository represents a repository for managing revisions.
type MigrationRepository struct {
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
func (x *MigrationRepository) LockRevision(ctx context.Context, params *LockRevisionParams) error {
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
func (x *MigrationRepository) UnlockRevision(ctx context.Context, params *UnlockRevisionParams) error {
	args := &ExecDeleteRevisionParams{}
	args.ID = params.Revision.ID
	return x.Gateway.ExecDeleteRevision(ctx, args)
}

// ApplyMigrationParams represents the parameters for executing a revision.
type ApplyMigrationParams struct {
	// Migration contains the parameters for executing a migration.
	Migration *Migration
}

// ApplyMigration executes a revision.
func (x *MigrationRepository) ApplyMigration(ctx context.Context, params *ApplyMigrationParams) error {
	repository := &JobRepository{
		Gateway: x.Gateway,
	}

	args := &UpsertRevisionParams{}
	args.SetRevision(params.Migration.Revision)
	// prepare the revision
	revision, err := x.Gateway.UpsertRevision(ctx, args)
	if err != nil {
		return err
	}

	params.Migration.Revision = revision

	start := time.Now()
	// Apply the statements one by one
	for index, query := range params.Migration.Statements {
		if index+1 <= params.Migration.Revision.Count {
			continue
		}

		query = commentRegexp.ReplaceAllString(query, "")
		query = concurrentlyRegexp.ReplaceAllString(query, "ASYNC")
		query = strings.TrimSpace(query)

		if len(query) > 0 {
			// execute the revision
			row := x.Gateway.Database().QueryRow(ctx, query)
			// Some queries returns a job id because they are asynchronous
			var jid string
			err := row.Scan(&jid)

			switch {
			case err == pgx.ErrNoRows:
				// We are good to go because the operation does not return job id
			case err != nil:
				msg := err.Error()
				// set the revision error
				params.Migration.Revision.Error = &msg
			default:
				args := &WaitJobParams{}
				args.JobID = jid
				// Wait for the job to complete
				job, err := repository.WaitJob(ctx, args)
				switch {
				case err == pgx.ErrNoRows:
				case err != nil:
					msg := err.Error()
					// set the revision error
					params.Migration.Revision.Error = &msg
				case job.Status == "failed":
					// set the revision error
					params.Migration.Revision.Error = job.Details
				}
			}
		}

		params.Migration.Revision.Count = index + 1
		params.Migration.Revision.ErrorStmt = &query
		params.Migration.Revision.ExecutedAt = time.Now().UTC()
		params.Migration.Revision.ExecutionTime = time.Since(start)

		args := &ExecUpdateRevisionParams{}
		args.SetRevision(params.Migration.Revision)

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

		if args.Error != nil {
			return nil
		}
	}

	return nil
}

// ListMigrationsParams represents the parameters for listing migrations.
type ListMigrationsParams struct{}

// ListMigrations lists all revisions in the repository.
func (x *MigrationRepository) ListMigrations(ctx context.Context, _ *ListMigrationsParams) (collection []*Migration, _ error) {
	matches, err := fs.Glob(x.FileSystem, "*.sql")
	if err != nil {
		return nil, err
	}

	fmt.Printf("Found %d migration files\n", len(matches))

	for _, path := range sort.StringSlice(matches) {
		// read the revision content
		data, err := fs.ReadFile(x.FileSystem, path)
		if err != nil {
			return nil, err
		}

		migration := &Migration{}
		migration.Statements = strings.SplitAfter(string(data), ";")
		migration.Revision = &Revision{}
		migration.Revision.SetName(path)
		migration.Revision.Total = len(migration.Statements)

		// append the revision
		collection = append(collection, migration)

		params := &GetRevisionParams{}
		params.SetRevision(migration.Revision)

		// load the revision
		revision, err := x.Gateway.GetRevision(ctx, params)
		switch {
		case err == pgx.ErrNoRows:
		// We are good to go because the revision does not exist
		case err != nil:
			return nil, err
		default:
			migration.Revision = revision
		}

		collection = append(collection, migration)
	}

	return collection, nil
}
