package sc

import (
	"fmt"
	"github.com/turbolytics/collector/internal/collector/state"
	"github.com/turbolytics/collector/internal/collector/state/memory"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
)

type StateStore struct {
	Type   state.StoreType
	Storer state.Storer
	Config map[string]any
}

type Config struct {
	StateStore StateStore `yaml:"state_store"`

	logger *zap.Logger
}

type Option func(*Config)

func WithLogger(l *zap.Logger) Option {
	return func(c *Config) {
		c.logger = l
	}
}

func NewFromFile(fname string, opts ...Option) (*Config, error) {
	bs, err := os.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	return New(bs, opts...)
}

func New(raw []byte, opts ...Option) (*Config, error) {
	var conf Config

	for _, opt := range opts {
		opt(&conf)
	}

	if err := yaml.Unmarshal(raw, &conf); err != nil {
		return nil, err
	}

	if err := initState(&conf); err != nil {
		return nil, err
	}

	return &conf, nil
}

func initState(c *Config) error {
	var s state.Storer
	var err error
	switch c.StateStore.Type {
	case state.StoreTypeMemory:
		s, err = memory.NewFromGenericConfig(c.StateStore.Config)
	default:
		return fmt.Errorf("state store type: %q unknown", c.StateStore.Type)
	}

	if err != nil {
		return err
	}
	c.StateStore.Storer = s
	return nil
}
