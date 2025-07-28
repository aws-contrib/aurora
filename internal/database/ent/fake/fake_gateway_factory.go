package fake

// NewFakeGateway returns a new instance of FakeGateway with a predefined transaction return value.
func NewFakeGateway() *FakeGateway {
	gateway := &FakeGateway{}
	gateway.DatabaseReturns(NewFakeDBTX())
	// revision operations
	gateway.GetRevisionReturns(NewFakeRevision(), nil)
	gateway.InsertRevisionReturns(NewFakeRevision(), nil)
	gateway.UpsertRevisionReturns(NewFakeRevision(), nil)
	gateway.UpdateRevisionReturns(NewFakeRevision(), nil)
	gateway.DeleteRevisionReturns(NewFakeRevision(), nil)
	// job operations
	gateway.GetJobReturns(NewFakeJob(), nil)
	gateway.InsertJobReturns(NewFakeJob(), nil)
	gateway.DeleteJobReturns(NewFakeJob(), nil)
	// lock operations
	gateway.GetLockReturns(NewFakeLock(), nil)
	gateway.InsertLockReturns(NewFakeLock(), nil)
	gateway.DeleteLockReturns(NewFakeLock(), nil)

	return gateway
}
