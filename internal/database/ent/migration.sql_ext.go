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

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var (
	commentRegexp      = regexp.MustCompile(`(?m)^\s*--.*$`)
	indexRegexp        = regexp.MustCompile(`(?i)INDEX`)
	concurrentlyRegexp = regexp.MustCompile(`(?i)CONCURRENTLY`)
)

//counterfeiter:generate -o ./fake . FileSystem

// FileSystem represents a filesystem that supports globbing and reading files.
type FileSystem interface {
	fs.FS
	fs.GlobFS
	fs.ReadFileFS
}

// MigrationLock is a UUID used to identify the migration lock in the database.
var MigrationLock = uuid.NewMD5(uuid.NameSpaceOID, []byte("aurora_schema_migrations"))

// MigrationRepository represents a repository for managing revisions.
type MigrationRepository struct {
	// Gateway represents the database gateway.
	Gateway Gateway
	// FileSystem is the filesystem where the revision files are located.
	FileSystem fs.FS
}

// LockMigrationParams represents the parameters for locking a revision.
type LockMigrationParams struct {
	// Timeout is the maximum time to wait for the lock.
	Timeout time.Duration
}

// LockMigration locks a revision for exclusive access.
func (x *MigrationRepository) LockMigration(ctx context.Context, params *LockMigrationParams) error {
	start := time.Now()

	for {
		// create the revision
		args := &ExecInsertLockParams{}
		args.ID = MigrationLock.String()
		args.CreatedAt = time.Now().UTC()

		err := x.Gateway.ExecInsertLock(ctx, args)
		// Waiting for the lock to be acquired
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

// UnlockMigration unlocks the revision after exclusive access.
func (x *MigrationRepository) UnlockMigration(ctx context.Context) error {
	args := &ExecDeleteLockParams{}
	args.ID = MigrationLock.String()

	return x.Gateway.ExecDeleteLock(ctx, args)
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

	start := time.Now()
	// Apply the statements one by one
	for index, query := range params.Migration.Statements {
		if index+1 <= revision.Count {
			continue
		}

		query = commentRegexp.ReplaceAllString(query, "")
		query = concurrentlyRegexp.ReplaceAllString(query, "")
		query = indexRegexp.ReplaceAllString(query, "INDEX ASYNC")
		query = strings.TrimSpace(query)

		if len(query) > 0 {
			// execute the revision
			row := x.Gateway.Database().QueryRow(ctx, query)
			// Some queries returns a job id because they are asynchronous
			var jid string
			err := row.Scan(&jid)

			switch {
			case err == pgx.ErrNoRows:
				revision.Count = index + 1
			case err != nil:
				msg := err.Error()
				// set the revision error
				revision.Error = &msg
				revision.ErrorStmt = &query
			default:
				args := &WaitJobParams{}
				args.JobID = jid
				// Wait for the job to complete
				job, err := repository.WaitJob(ctx, args)
				switch {
				case err == pgx.ErrNoRows:
					revision.Count = index + 1
				case err != nil:
					msg := err.Error()
					// set the revision error
					revision.Error = &msg
					revision.ErrorStmt = &query
				case job.Status == "failed":
					// set the revision error
					revision.Error = job.Details
					revision.ErrorStmt = &query
				default:
					revision.Count = index + 1
				}
			}
		} else {
			revision.Count = index + 1
		}

		revision.ExecutedAt = time.Now().UTC()
		revision.ExecutionTime = time.Since(start)

		args := &ExecUpdateRevisionParams{}
		args.SetRevision(revision)
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

		// Update the migration parameters
		params.Migration.Revision = revision
		// Stop processing if there is an error
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
