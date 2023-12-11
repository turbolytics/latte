package http

import "github.com/mitchellh/mapstructure"

type config struct {
	server string
	path   string
}

type HTTP struct{}

func (h *HTTP) Write(bs []byte) (int, error) {
	return 0, nil
}

func NewFromGenericConfig(m map[string]any) (*HTTP, error) {
	var conf config
	if err := mapstructure.Decode(m, &conf); err != nil {
		return nil, err
	}

	return &HTTP{}, nil
}
