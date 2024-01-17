package service

import (
	"context"
	"github.com/go-co-op/gocron/v2"
	"github.com/turbolytics/collector/internal/collector"
	"go.uber.org/zap"
)

type Service struct {
	logger     *zap.Logger
	collectors []collector.Collector
	scheduler  gocron.Scheduler
}

func (s *Service) Shutdown() {
	s.logger.Info("shutdown")
	defer s.scheduler.Shutdown()
	// close each collector
	for _, coll := range s.collectors {
		coll.Close()
	}

}

func (s *Service) Run(ctx context.Context) error {
	s.logger.Info("run")
	// iterate all collectors and invoke for the initial invocation
	for _, col := range s.collectors {
		go col.InvokeHandleError(ctx)
	}

	for _, col := range s.collectors {
		colCopy := col

		var jd gocron.JobDefinition
		if colCopy.Interval() != nil {
			jd = gocron.DurationJob(
				*(colCopy.Interval()),
			)
		} else if colCopy.Cron() != nil {
			jd = gocron.CronJob(
				*(colCopy.Cron()),
				false,
			)
		}

		_, err := s.scheduler.NewJob(
			jd,
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

func NewService(cs []collector.Collector, opts ...Option) (*Service, error) {
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
