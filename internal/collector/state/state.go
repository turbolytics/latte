package state

import (
	"context"
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
