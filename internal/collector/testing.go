package collector

type TestSink struct {
	Closes int
}

func (ts *TestSink) Write(bs []byte) (int, error) {
	return 0, nil
}

func (ts *TestSink) Close() error {
	ts.Closes++
	return nil
}
