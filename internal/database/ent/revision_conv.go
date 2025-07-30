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
