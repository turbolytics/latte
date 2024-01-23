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

type Storer interface {
	io.Closer

	MostRecentInvocation(ctx context.Context, collectorName string) (*Invocation, error)
	SaveInvocation(invocation *Invocation) error
}

type Store struct {
	Type   StoreType
	Storer Storer
	Config map[string]any
}

func (s Store) Validate() error {
	ts := map[StoreType]struct{}{
		"":              {},
		StoreTypeMemory: {},
	}

	if _, ok := ts[s.Type]; !ok {
		fmt.Errorf("unknown strategy: %v", s.Type)
	}

	return nil
}
