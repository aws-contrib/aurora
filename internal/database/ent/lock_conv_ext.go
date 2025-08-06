//go:build !goverter

package ent

// SetLock sets the params from the entity.
func (x *GetLockParams) SetLock(entity *Lock) {
	converter := &GetLockParamsConverterImpl{}
	converter.SetFromLock(x, entity)
}

// SetLock sets the params from the entity.
func (x *InsertLockParams) SetLock(entity *Lock) {
	converter := &InsertLockParamsConverterImpl{}
	converter.SetFromLock(x, entity)
}

// SetLock sets the params from the entity.
func (x *ExecInsertLockParams) SetLock(entity *Lock) {
	converter := &ExecInsertLockParamsConverterImpl{}
	converter.SetFromLock(x, entity)
}

// SetLock sets the params from the entity.
func (x *DeleteLockParams) SetLock(entity *Lock) {
	converter := &DeleteLockParamsConverterImpl{}
	converter.SetFromLock(x, entity)
}

// SetLock sets the params from the entity.
func (x *ExecDeleteLockParams) SetLock(entity *Lock) {
	converter := &ExecDeleteLockParamsConverterImpl{}
	converter.SetFromLock(x, entity)
}
