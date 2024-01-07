package collector

type TestSink struct {
	closes int
}

func (ts *TestSink) Write(bs []byte) (int, error) {
	return 0, nil
}

func (ts *TestSink) Close() error {
	ts.closes++
	return nil
}
