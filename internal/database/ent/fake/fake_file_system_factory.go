package fake

// NewFakeFileSystem returns a new instance of FakeFileSystem with a predefined FileSystem return value.
func NewFakeFileSystem() *FakeFileSystem {
	fs := &FakeFileSystem{}
	fs.ReadFileReturns([]byte("CREATE TABLE IF NOT EXISTS example (id SERIAL PRIMARY KEY, name VARCHAR(255));"), nil)
	fs.GlobReturns([]string{"aurora_schema_table_test.sql"}, nil)

	return fs
}
