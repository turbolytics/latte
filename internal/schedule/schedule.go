package schedule

import (
	"fmt"
	"time"
)

type Config struct {
	Interval *time.Duration
	Cron     *string
}

func (s Config) Validate() error {
	if s.Interval == nil && s.Cron == nil {
		return fmt.Errorf("must set invocation.interval or invocation.cron")
	}

	if s.Interval != nil && s.Cron != nil {
		return fmt.Errorf("must set either invocation.interval or invocation.cron")
	}

	return nil
}

type Schedule struct {
	config Config
}

func (s Schedule) Interval() *time.Duration {
	return s.config.Interval
}

func (s Schedule) Cron() *string {
	return s.config.Cron
}

func New(c Config) Schedule {
	return Schedule{
		config: c,
	}
}
