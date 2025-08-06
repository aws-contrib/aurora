package ent

import (
	"strings"
)

// Migration represents a database migration with its details.
type Migration struct {
	Revision   *Revision
	Statements []string
}

// MigrationState represents the state of a migration operation.
type MigrationState struct {
	Next     *Revision
	Current  *Revision
	Pending  []*Revision
	Executed []*Revision
}

// GetName returns the name of the revision file based on its ID and description.
func (x *Revision) GetName() string {
	return x.ID + "_" + x.Description + ".sql"
}

// SetName sets the name of the revision file.
func (x *Revision) SetName(name string) {
	name = strings.TrimSuffix(name, ".sql")

	for index, part := range strings.SplitN(name, "_", 2) {
		switch index {
		case 0:
			x.ID = part
		case 1:
			x.Description = part
		}
	}
}
