package config

type MockReader struct {
	OutputCfg   *Cfg
	OutputError error
}

func (mockReader *MockReader) ReadConfig(filename string) (*Cfg, error) {
	return mockReader.OutputCfg, mockReader.OutputError
}
