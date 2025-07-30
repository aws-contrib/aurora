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
