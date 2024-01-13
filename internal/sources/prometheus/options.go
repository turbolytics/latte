package prometheus

import (
	"github.com/turbolytics/collector/internal/collector/state"
	"github.com/turbolytics/collector/internal/invocation"
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

func WithInvocationStrategy(s invocation.TypeStrategy) Option {
	return func(p *Prometheus) {
		p.invocationStrategy = s
	}
}

func WithCollectorName(cn string) Option {
	return func(p *Prometheus) {
		p.collectorName = cn
	}
}
