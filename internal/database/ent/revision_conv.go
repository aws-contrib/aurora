package ent

// goverter:converter
// goverter:skipCopySameType yes
// goverter:output:file models_conv_gen.go
// goverter:output:package github.com/aws-contrib/aurora/internal/database/ent
type GetRevisionParamsConverter interface {
	// goverter:update target
	SetFromRevision(target *GetRevisionParams, source *Revision)
}

// goverter:converter
// goverter:skipCopySameType yes
// goverter:output:file models_conv_gen.go
// goverter:output:package github.com/aws-contrib/aurora/internal/database/ent
type InsertRevisionParamsConverter interface {
	// goverter:update target
	SetFromRevision(target *InsertRevisionParams, source *Revision)
}

// goverter:converter
// goverter:skipCopySameType yes
// goverter:output:file models_conv_gen.go
// goverter:output:package github.com/aws-contrib/aurora/internal/database/ent
type ExecInsertRevisionParamsConverter interface {
	// goverter:update target
	SetFromRevision(target *ExecInsertRevisionParams, source *Revision)
}

// goverter:converter
// goverter:skipCopySameType yes
// goverter:output:file models_conv_gen.go
// goverter:output:package github.com/aws-contrib/aurora/internal/database/ent
type UpsertRevisionParamsConverter interface {
	// goverter:update target
	SetFromRevision(target *UpsertRevisionParams, source *Revision)
}

// goverter:converter
// goverter:skipCopySameType yes
// goverter:output:file models_conv_gen.go
// goverter:output:package github.com/aws-contrib/aurora/internal/database/ent
type ExecUpsertRevisionParamsConverter interface {
	// goverter:update target
	SetFromRevision(target *ExecUpsertRevisionParams, source *Revision)
}

// goverter:converter
// goverter:skipCopySameType yes
// goverter:output:file models_conv_gen.go
// goverter:output:package github.com/aws-contrib/aurora/internal/database/ent
type UpdateRevisionParamsConverter interface {
	// goverter:update target
	// goverter:ignore UpdateMask
	SetFromRevision(target *UpdateRevisionParams, source *Revision)
}

// goverter:converter
// goverter:skipCopySameType yes
// goverter:output:file models_conv_gen.go
// goverter:output:package github.com/aws-contrib/aurora/internal/database/ent
type ExecUpdateRevisionParamsConverter interface {
	// goverter:update target
	// goverter:ignore UpdateMask
	SetFromRevision(target *ExecUpdateRevisionParams, source *Revision)
}

// goverter:converter
// goverter:skipCopySameType yes
// goverter:output:file models_conv_gen.go
// goverter:output:package github.com/aws-contrib/aurora/internal/database/ent
type DeleteRevisionParamsConverter interface {
	// goverter:update target
	SetFromRevision(target *DeleteRevisionParams, source *Revision)
}

// goverter:converter
// goverter:skipCopySameType yes
// goverter:output:file models_conv_gen.go
// goverter:output:package github.com/aws-contrib/aurora/internal/database/ent
type ExecDeleteRevisionParamsConverter interface {
	// goverter:update target
	SetFromRevision(target *ExecDeleteRevisionParams, source *Revision)
}
