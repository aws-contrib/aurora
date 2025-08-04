package fake

// NewFakeDBTX returns a new instance of FakeDBTX with a predefined QueryRow return value.
func NewFakeDBTX() *FakeDBTX {
	dbx := &FakeDBTX{}
	dbx.QueryRowReturns(&FakeRow{})
	return dbx
}
