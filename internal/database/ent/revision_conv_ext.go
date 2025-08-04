//go:build !goverter

package ent

// SetRevision sets the params from the entity.
func (x *GetRevisionParams) SetRevision(entity *Revision) {
	converter := &GetRevisionParamsConverterImpl{}
	converter.SetFromRevision(x, entity)
}

// SetRevision sets the params from the entity.
func (x *InsertRevisionParams) SetRevision(entity *Revision) {
	converter := &InsertRevisionParamsConverterImpl{}
	converter.SetFromRevision(x, entity)
}

// SetRevision sets the params from the entity.
func (x *ExecInsertRevisionParams) SetRevision(entity *Revision) {
	converter := &ExecInsertRevisionParamsConverterImpl{}
	converter.SetFromRevision(x, entity)
}

// SetRevision sets the params from the entity.
func (x *UpsertRevisionParams) SetRevision(entity *Revision) {
	converter := &UpsertRevisionParamsConverterImpl{}
	converter.SetFromRevision(x, entity)
}

// SetRevision sets the params from the entity.
func (x *ExecUpsertRevisionParams) SetRevision(entity *Revision) {
	converter := &ExecUpsertRevisionParamsConverterImpl{}
	converter.SetFromRevision(x, entity)
}

// SetRevision sets the params from the entity.
func (x *UpdateRevisionParams) SetRevision(entity *Revision) {
	converter := &UpdateRevisionParamsConverterImpl{}
	converter.SetFromRevision(x, entity)
}

// SetRevision sets the params from the entity.
func (x *ExecUpdateRevisionParams) SetRevision(entity *Revision) {
	converter := &ExecUpdateRevisionParamsConverterImpl{}
	converter.SetFromRevision(x, entity)
}

// SetRevision sets the params from the entity.
func (x *DeleteRevisionParams) SetRevision(entity *Revision) {
	converter := &DeleteRevisionParamsConverterImpl{}
	converter.SetFromRevision(x, entity)
}

// SetRevision sets the params from the entity.
func (x *ExecDeleteRevisionParams) SetRevision(entity *Revision) {
	converter := &ExecDeleteRevisionParamsConverterImpl{}
	converter.SetFromRevision(x, entity)
}
