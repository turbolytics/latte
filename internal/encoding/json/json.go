package json

import (
	"bytes"
	"encoding/json"
	"io"
)

type JSON struct {
	buf io.Writer
}

func (j *JSON) Flush() error {
	return nil
}

func (j *JSON) Close() error {
	return nil
}

func (j *JSON) Init(buf *bytes.Buffer) error {
	j.buf = buf
	return nil
}

func (j *JSON) Write(d any) error {
	bs, err := json.Marshal(d)
	if err != nil {
		return err
	}
	if _, err = j.buf.Write(bs); err != nil {
		return err
	}

	_, err = j.buf.Write([]byte("\n"))

	return err
}

func NewFromGenericConfig(m map[string]any) (*JSON, error) {
	j := &JSON{}
	return j, nil
}
