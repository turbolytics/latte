package collector

import (
	"context"
	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
)

type Service struct {
	logger     *zap.Logger
	collectors []*Collector
	scheduler  gocron.Scheduler
}

func (s *Service) Shutdown() {
	s.logger.Info("shutdown")
	s.scheduler.Shutdown()
}

func (s *Service) Run(ctx context.Context) error {
	s.logger.Info("run")
	// iterate all collectors and invoke for the initial invocation
	for _, col := range s.collectors {
		go col.InvokeHandleError(ctx)
	}

	for _, col := range s.collectors {
		_, err := s.scheduler.NewJob(
			gocron.DurationJob(
				col.Config.Schedule.Interval,
			),
			gocron.NewTask(
				func(ctx context.Context) {
					col.InvokeHandleError(ctx)
				},
				ctx,
			),
		)
		if err != nil {
			return err
		}
	}

	s.scheduler.Start()
	<-ctx.Done()

	return nil
}

type ServiceOption func(*Service)

func WithServiceLogger(l *zap.Logger) ServiceOption {
	return func(s *Service) {
		s.logger = l
	}
}

func NewService(cs []*Collector, opts ...ServiceOption) (*Service, error) {
	sch, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	s := &Service{
		collectors: cs,
		scheduler:  sch,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s, nil
}
