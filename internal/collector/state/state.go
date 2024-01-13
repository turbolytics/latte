package state

import (
	"time"
)

type StoreType string

const (
	StoreTypeMemory StoreType = "memory"
)

type Invocation struct {
	CollectorName  string
	WindowDatetime time.Time
	Time           time.Time
}

type Storer interface {
	MostRecentInvocation(collectorName string) (*Invocation, error)
	SaveInvocation(invocation *Invocation) error
}
