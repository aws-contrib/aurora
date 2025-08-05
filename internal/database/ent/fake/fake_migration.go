package fake

import "github.com/aws-contrib/aurora/internal/database/ent"

// NewFakeMigration creates a new fake migration entity with a revision.
func NewFakeMigration() *ent.Migration {
	entity := &ent.Migration{}
	entity.Revision = NewFakeRevision()
	entity.Statements = []string{
		"CREATE TABLE IF NOT EXISTS example (id SERIAL PRIMARY KEY, name VARCHAR(255));",
	}

	return entity
}
