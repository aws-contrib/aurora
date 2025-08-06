package fake

import (
	"time"

	"github.com/aws-contrib/aurora/internal/database/ent"
	"github.com/google/uuid"
)

// NewFakeLock returns a new fake revision.
func NewFakeLock() *ent.Lock {
	return &ent.Lock{
		ID:        uuid.New().String(),
		CreatedAt: time.Now().Truncate(time.Millisecond),
	}
}
