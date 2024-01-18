package collector

import (
	"context"
	"time"
)

type Collector interface {
	Close() error
	InvokeHandleError(context.Context)
	Interval() *time.Duration
	Cron() *string
}
