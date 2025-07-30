package fake

import (
	"time"

	"github.com/aws-contrib/aurora/internal/database/ent"
	"github.com/google/uuid"
)

// NewFakeRevision returns a new fake revision.
func NewFakeRevision() *ent.Revision {
	return &ent.Revision{
		ID:            uuid.New(),
		Description:   "schema",
		ExecutedAt:    time.Now().Truncate(time.Millisecond),
		ExecutionTime: 500 * time.Millisecond,
	}
}
