package source

import (
	"context"
	"time"
)

type Record interface {
	Bytes() ([]byte, error)
}

type Result interface {
	Transform() error
	Records() []Record
}

type Sourcer interface {
	Source(ctx context.Context) (Result, error)
	Window() *time.Duration
}
