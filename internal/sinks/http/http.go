package http

import (
	"bytes"
	"github.com/mitchellh/mapstructure"
	"github.com/turbolytics/latte/internal/sinks"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type config struct {
	URI     string
	Method  string
	Body    string
	Headers map[string]string
}

type Option func(*HTTP)

func WithLogger(l *zap.Logger) Option {
	return func(h *HTTP) {
		h.logger = l
	}
}

type HTTP struct {
	config config

	logger *zap.Logger
}

func (h *HTTP) Close() error {
	return nil
}

func (h *HTTP) Type() sinks.Type {
	return sinks.TypeHTTP
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
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	h.logger.Debug(
		"http.response",
		zap.String("name", "http.sink"),
		zap.ByteString("resp", body),
	)

	return len(bs), nil
}

func NewFromGenericConfig(m map[string]any, opts ...Option) (*HTTP, error) {
	var conf config
	if err := mapstructure.Decode(m, &conf); err != nil {
		return nil, err
	}

	h := &HTTP{
		config: conf,
	}
	for _, opt := range opts {
		opt(h)
	}

	return h, nil
}
