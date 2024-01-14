package state

import (
	"github.com/turbolytics/collector/internal/timeseries"
	"time"
)

type StoreType string

const (
	StoreTypeMemory StoreType = "memory"
)

type Invocation struct {
	CollectorName string
	Time          time.Time
	Window        *timeseries.Bucket
}

type Storer interface {
	MostRecentInvocation(collectorName string) (*Invocation, error)
	SaveInvocation(invocation *Invocation) error
}
