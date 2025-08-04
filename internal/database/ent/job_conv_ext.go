//go:build !goverter

package ent

// SetJob sets the params from the entity.
func (x *GetJobParams) SetJob(entity *Job) {
	converter := &GetJobParamsConverterImpl{}
	converter.SetFromJob(x, entity)
}

// SetJob sets the params from the entity.
func (x *InsertJobParams) SetJob(entity *Job) {
	converter := &InsertJobParamsConverterImpl{}
	converter.SetFromJob(x, entity)
}

// SetJob sets the params from the entity.
func (x *ExecInsertJobParams) SetJob(entity *Job) {
	converter := &ExecInsertJobParamsConverterImpl{}
	converter.SetFromJob(x, entity)
}

// SetJob sets the params from the entity.
func (x *DeleteJobParams) SetJob(entity *Job) {
	converter := &DeleteJobParamsConverterImpl{}
	converter.SetFromJob(x, entity)
}

// SetJob sets the params from the entity.
func (x *ExecDeleteJobParams) SetJob(entity *Job) {
	converter := &ExecDeleteJobParamsConverterImpl{}
	converter.SetFromJob(x, entity)
}
