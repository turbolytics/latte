package state

import (
	"context"
	"fmt"
	"github.com/turbolytics/collector/internal/timeseries"
	"io"
	"time"
)

type StoreType string

const (
	StoreTypeMemory StoreType = "memory"
)

type Invocation struct {
	CollectorName string
	Time          time.Time
	Window        *timeseries.Window
}

func (i Invocation) End() *time.Time {
	if i.Window != nil {
		return &i.Window.End
	}
	return nil
}

type Storer interface {
	io.Closer

	MostRecentInvocation(ctx context.Context, collectorName string) (*Invocation, error)
	SaveInvocation(invocation *Invocation) error
}

type Config struct {
	Type   StoreType
	Storer Storer
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

func (c *Config) Init() error {
	var s Storer
	var err error
	switch c.Type {
	case StoreTypeMemory:
		s, err = NewMemoryStoreFromGenericConfig(c.Config)
	}
	if err != nil {
		return err
	}
	c.Storer = s
	return nil
}
