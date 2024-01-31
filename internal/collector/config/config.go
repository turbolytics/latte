package config

import (
	"github.com/turbolytics/latte/internal/collector/schedule"
	collSource "github.com/turbolytics/latte/internal/collector/source"
	"github.com/turbolytics/latte/internal/sinks"
)

type Config interface {
	CollectorName() string
	GetSinks() []sinks.Sinker
	GetSchedule() schedule.Schedule
	GetSource() collSource.Config
}
