package ent

// goverter:converter
// goverter:skipCopySameType yes
// goverter:output:file models_conv_gen.go
// goverter:output:package github.com/aws-contrib/aurora/internal/database/ent
type GetLockParamsConverter interface {
	// goverter:update target
	SetFromLock(target *GetLockParams, source *Lock)
}

// goverter:converter
// goverter:skipCopySameType yes
// goverter:output:file models_conv_gen.go
// goverter:output:package github.com/aws-contrib/aurora/internal/database/ent
type InsertLockParamsConverter interface {
	// goverter:update target
	SetFromLock(target *InsertLockParams, source *Lock)
}

// goverter:converter
// goverter:skipCopySameType yes
// goverter:output:file models_conv_gen.go
// goverter:output:package github.com/aws-contrib/aurora/internal/database/ent
type ExecInsertLockParamsConverter interface {
	// goverter:update target
	SetFromLock(target *ExecInsertLockParams, source *Lock)
}

// goverter:converter
// goverter:skipCopySameType yes
// goverter:output:file models_conv_gen.go
// goverter:output:package github.com/aws-contrib/aurora/internal/database/ent
type DeleteLockParamsConverter interface {
	// goverter:update target
	SetFromLock(target *DeleteLockParams, source *Lock)
}

// goverter:converter
// goverter:skipCopySameType yes
// goverter:output:file models_conv_gen.go
// goverter:output:package github.com/aws-contrib/aurora/internal/database/ent
type ExecDeleteLockParamsConverter interface {
	// goverter:update target
	SetFromLock(target *ExecDeleteLockParams, source *Lock)
}
