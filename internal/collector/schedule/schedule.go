package schedule

import (
	"fmt"
	"time"
)

type Schedule struct {
	Interval *time.Duration
	Cron     *string
}

func (s Schedule) Validate() error {
	if s.Interval == nil && s.Cron == nil {
		return fmt.Errorf("must set invocation.interval or invocation.cron")
	}

	if s.Interval != nil && s.Cron != nil {
		return fmt.Errorf("must set either invocation.interval or invocation.cron")
	}

	return nil
}
