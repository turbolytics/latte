package http

import (
	"bytes"
	"github.com/mitchellh/mapstructure"
	"io"
	"net/http"
)

type config struct {
	URI     string
	Method  string
	Body    string
	Headers map[string]string
}

type HTTP struct {
	config config
}

func (h *HTTP) Close() error {
	return nil
}

func (h *HTTP) Write(bs []byte) (int, error) {
	buf := bytes.NewBuffer(bs)

	req, err := http.NewRequest(
		h.config.Method,
		h.config.URI,
		buf,
	)
	if err != nil {
		return 0, err
	}

	for k, v := range h.config.Headers {
		req.Header.Add(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	return len(bs), nil
}

func NewFromGenericConfig(m map[string]any) (*HTTP, error) {
	var conf config
	if err := mapstructure.Decode(m, &conf); err != nil {
		return nil, err
	}

	return &HTTP{
		config: conf,
	}, nil
}
