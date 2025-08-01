//go:build !goverter

package ent

import (
	"context"
	"fmt"
	"io/fs"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

var asyncre = regexp.MustCompile(`(?i)CONCURRENTLY`)

//counterfeiter:generate -o ./fake io/fs.ReadFileFS

type Applier interface {
	ApplyRevision(ctx context.Context, params *ApplyRevisionParams) error
}

// ApplyRevisionParams represents the parameters for executing a revision.
type ApplyRevisionParams struct {
	// Revision contains the parameters for executing a revision.
	Revision *Revision
	// FileSystem is the filesystem where the revision files are located.
	FileSystem fs.FS
}

// ApplyRevision executes a revision.
func (x *Queries) ApplyRevision(ctx context.Context, params *ApplyRevisionParams) error {
	data, err := fs.ReadFile(params.FileSystem, params.Revision.GetName())
	if err != nil {
		return err
	}

	start := time.Now()
	// Apply the statements one by one
	for _, query := range strings.SplitAfter(string(data), ";") {
		query = asyncre.ReplaceAllString(query, "ASYNC")

		params := &WaitForJobParams{}
		// execute the revision
		rerr := x.db.QueryRow(ctx, query).Scan(&params.JobID)

		switch {
		case rerr == pgx.ErrNoRows:
		// We are good to go because the operation does not return job id
		case rerr != nil:
			return rerr
		default:
			ok, werr := x.WaitForJob(ctx, params)
			switch {
			case werr == pgx.ErrNoRows:
				// We are good to go because the operation does not return job id
			case werr != nil:
				return werr
			case !ok:
				return fmt.Errorf("job %s did not complete successfully", params.JobID)
			}
		}
	}

	// prepare the parameters
	params.Revision.ExecutedAt = time.Now().UTC()
	params.Revision.ExecutionTime = time.Since(start)
	// create the revision
	args := &InsertRevisionParams{}
	args.SetRevision(params.Revision)
	_, err = x.InsertRevision(ctx, args)
	return err
}
