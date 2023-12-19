package service

import (
	"context"
	"github.com/go-co-op/gocron/v2"
	"github.com/turbolytics/collector/internal/collector"
	"go.uber.org/zap"
)

type Service struct {
	logger     *zap.Logger
	collectors []*collector.Collector
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
		colCopy := col
		_, err := s.scheduler.NewJob(
			gocron.DurationJob(
				colCopy.Config.Schedule.Interval,
			),
			gocron.NewTask(
				func(ctx context.Context) {
					colCopy.InvokeHandleError(ctx)
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

type Option func(*Service)

func WithLogger(l *zap.Logger) Option {
	return func(s *Service) {
		s.logger = l
	}
}

func NewService(cs []*collector.Collector, opts ...Option) (*Service, error) {
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