package fake

import (
	"github.com/aws-contrib/aurora/internal/database/ent"
	"github.com/google/uuid"
)

// NewFakeJob returns a new job
func NewFakeJob() *ent.Job {
	return &ent.Job{
		ID:      uuid.New().String(),
		Status:  "completed",
		Details: "",
	}
}
