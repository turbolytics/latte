package prometheus

import (
	"github.com/turbolytics/collector/internal/collector/state"
	"github.com/turbolytics/collector/internal/schedule"
	"go.uber.org/zap"
)

type Option func(*Prometheus)

func WithLogger(l *zap.Logger) Option {
	return func(p *Prometheus) {
		p.logger = l
	}
}

func WithStateStorer(ss state.Storer) Option {
	return func(p *Prometheus) {
		p.stateStorer = ss
	}
}

func WithScheduleStrategy(ss schedule.TypeStrategy) Option {
	return func(p *Prometheus) {
		p.scheduleStrategy = ss
	}
}
