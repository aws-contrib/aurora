//go:build !goverter

package ent

import (
	"context"
	"io/fs"
	"strings"
	"time"
)

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
		// execute the revision
		if _, err = x.db.Exec(ctx, query); err != nil {
			return err
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
