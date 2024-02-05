package service

import (
	"context"
	"github.com/go-co-op/gocron/v2"
	"github.com/turbolytics/latte/internal/invoker"
	"go.uber.org/zap"
)

type Service struct {
	logger    *zap.Logger
	invokers  []*invoker.Invoker
	scheduler gocron.Scheduler
}

func (s *Service) Shutdown() {
	s.logger.Info("shutdown")
	defer s.scheduler.Shutdown()
	// close each collector
	for _, i := range s.invokers {
		i.Close()
	}

}

func (s *Service) Run(ctx context.Context) error {
	s.logger.Info("run")
	// iterate all collectors and invoke for the initial invocation
	for _, i := range s.invokers {
		go i.InvokeHandleError(ctx)
	}

	for _, i := range s.invokers {
		iCopy := i
		sch := i.Collector.GetSchedule()

		var jd gocron.JobDefinition
		if sch.Interval != nil {
			jd = gocron.DurationJob(
				*(sch.Interval),
			)
		} else if sch.Cron != nil {
			jd = gocron.CronJob(
				*(sch.Cron),
				false,
			)
		}

		_, err := s.scheduler.NewJob(
			jd,
			gocron.NewTask(
				func(ctx context.Context) {
					iCopy.InvokeHandleError(ctx)
				},
				ctx,
			),
			gocron.WithSingletonMode(gocron.LimitModeReschedule),
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

func NewService(is []*invoker.Invoker, opts ...Option) (*Service, error) {
	sch, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	s := &Service{
		invokers:  is,
		scheduler: sch,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s, nil
}
