package state

import (
	"time"
)

type StoreType string

const (
	StoreTypeMemory StoreType = "memory"
)

type Invocation struct {
	Collector     string
	GrainDatetime time.Time
	Time          time.Time
}

type Storer interface {
	LastInvocation(collector string) (*Invocation, error)
	SaveInvocation(invocation *Invocation) error
}
