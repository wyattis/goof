package goof

type testModule struct {
	ModuleBase
}

func (m *testModule) Id() string {
	return "test"
}

func init() {
	var _ Module = &testModule{}
}
