package state

import (
	"fmt"
)

type StoreType string

const (
	StoreTypeMemory StoreType = "memory"
)

type Config struct {
	Type   StoreType
	Config map[string]any
}

func (c Config) Validate() error {
	ts := map[StoreType]struct{}{
		"":              {},
		StoreTypeMemory: {},
	}

	if _, ok := ts[c.Type]; !ok {
		fmt.Errorf("unknown strategy: %v", c.Type)
	}

	return nil
}
