package config

import (
	"fmt"
	"github.com/turbolytics/collector/internal/schedule"
	"time"
)

type Schedule struct {
	Interval *time.Duration
	Cron     *string
	Strategy schedule.TypeStrategy
}

func (s Schedule) Validate() error {
	if s.Interval == nil && s.Cron == nil {
		return fmt.Errorf("must set schedule.interval or schedule.cron")
	}

	if s.Interval != nil && s.Cron != nil {
		return fmt.Errorf("must set either schedule.interval or schedule.cron")
	}

	vs := map[schedule.TypeStrategy]struct{}{
		schedule.TypeStrategyTick:     {},
		schedule.TypeStrategyStateful: {},
	}

	if _, ok := vs[s.Strategy]; !ok {
		return fmt.Errorf("unknown strategy: %q", s.Strategy)
	}

	return nil
}

func (s *Schedule) SetDefaults() {
	if s.Strategy == "" {
		s.Strategy = schedule.TypeStrategyTick
	}
}
