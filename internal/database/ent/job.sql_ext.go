//go:build !goverter

package ent

import (
	"context"
	"time"
)

// JobRepository provides methods to interact with the Job entity.
type JobRepository struct {
	Gateway Gateway
}

// WaitJobParams is the parameters for the WaitJob method.
type WaitJobParams struct {
	// Job is the job to wait for.
	JobID string
}

// WaitJob waits for a job to complete and returns the job details.
func (x *JobRepository) WaitJob(ctx context.Context, params *WaitJobParams) (*Job, error) {
	args := &GetJobParams{}
	args.JobID = params.JobID

	for {
		job, err := x.Gateway.GetJob(ctx, args)
		switch {
		case err != nil:
			return nil, err
		case
			job.Status == "submitted",
			job.Status == "processing":
			time.Sleep(100 * time.Millisecond)
		default:
			return job, nil
		}
	}
}
