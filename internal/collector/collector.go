package collector

import (
	"context"
	"time"
)

type Collector interface {
	Close() error
	InvokeHandleError(context.Context)
	Invoke(context.Context) error
	Interval() *time.Duration
	Cron() *string
}
