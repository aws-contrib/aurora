package ent

// goverter:converter
// goverter:skipCopySameType yes
// goverter:output:file models_conv_gen.go
// goverter:output:package github.com/aws-contrib/aurora/internal/database/ent
type GetJobParamsConverter interface {
	// goverter:update target
	SetFromJob(target *GetJobParams, source *Job)
}

// goverter:converter
// goverter:skipCopySameType yes
// goverter:output:file models_conv_gen.go
// goverter:output:package github.com/aws-contrib/aurora/internal/database/ent
type InsertJobParamsConverter interface {
	// goverter:update target
	SetFromJob(target *InsertJobParams, source *Job)
}

// goverter:converter
// goverter:skipCopySameType yes
// goverter:output:file models_conv_gen.go
// goverter:output:package github.com/aws-contrib/aurora/internal/database/ent
type ExecInsertJobParamsConverter interface {
	// goverter:update target
	SetFromJob(target *ExecInsertJobParams, source *Job)
}

// goverter:converter
// goverter:skipCopySameType yes
// goverter:output:file models_conv_gen.go
// goverter:output:package github.com/aws-contrib/aurora/internal/database/ent
type DeleteJobParamsConverter interface {
	// goverter:update target
	SetFromJob(target *DeleteJobParams, source *Job)
}

// goverter:converter
// goverter:skipCopySameType yes
// goverter:output:file models_conv_gen.go
// goverter:output:package github.com/aws-contrib/aurora/internal/database/ent
type ExecDeleteJobParamsConverter interface {
	// goverter:update target
	SetFromJob(target *ExecDeleteJobParams, source *Job)
}
