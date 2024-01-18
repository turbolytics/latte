package config

import (
	"fmt"
	"github.com/turbolytics/collector/internal/collector/state"
	"github.com/turbolytics/collector/internal/collector/state/memory"
)

type StateStore struct {
	Type   state.StoreType
	Storer state.Storer
	Config map[string]any
}

func (s StateStore) Validate() error {
	ts := map[state.StoreType]struct{}{
		"":                    {},
		state.StoreTypeMemory: {},
	}

	if _, ok := ts[s.Type]; !ok {
		fmt.Errorf("unknown strategy: %v", s.Type)
	}

	return nil
}

func initStateStore(c *Config) error {
	var s state.Storer
	var err error
	switch c.StateStore.Type {
	case state.StoreTypeMemory:
		s, err = memory.NewFromGenericConfig(c.StateStore.Config)
	}

	if err != nil {
		return err
	}
	c.StateStore.Storer = s
	return nil
}
