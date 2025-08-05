package fake

// NewFakeGateway returns a new instance of FakeGateway with a predefined transaction return value.
func NewFakeGateway() *FakeGateway {
	gateway := &FakeGateway{}
	gateway.DatabaseReturns(NewFakeDBTX())
	gateway.InsertRevisionReturns(NewFakeRevision(), nil)
	gateway.UpsertRevisionReturns(NewFakeRevision(), nil)
	gateway.UpdateRevisionReturns(NewFakeRevision(), nil)
	gateway.DeleteRevisionReturns(NewFakeRevision(), nil)
	gateway.GetRevisionReturns(NewFakeRevision(), nil)
	gateway.GetJobReturns(NewFakeJob(), nil)

	return gateway
}
