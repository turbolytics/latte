package http

import (
	"bytes"
	"github.com/mitchellh/mapstructure"
	"io"
	"net/http"
)

type config struct {
	URI string
}

type HTTP struct {
	config config
}

func (h *HTTP) Write(bs []byte) (int, error) {
	buf := bytes.NewBuffer(bs)
	resp, err := http.Post(h.config.URI, "application/x-www-form-urlencoded", buf)
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
