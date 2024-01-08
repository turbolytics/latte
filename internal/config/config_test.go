package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
	"time"
)

var exampleDir string

func init() {
	currDir, _ := os.Getwd()
	exampleDir = path.Join(currDir, "..", "..", "dev", "examples")
}

func TestNewConfigFromFile(t *testing.T) {
	testCases := []struct {
		fileName string
	}{
		{"postgres.http.stdout.yaml"},
		{"postgres.stdout.yaml"},
		{"mongo.http.stdout.yaml"},
		{"postgres.kafka.stdout.yaml"},
		{"postgres.fileaudit.stdout.yaml"},
		{"prometheus.stdout.yaml"},
	}
	for _, tc := range testCases {
		t.Run(tc.fileName, func(t *testing.T) {
			fPath := path.Join(exampleDir, tc.fileName)
			_, err := NewFromFile(
				fPath,
				WithJustValidation(true),
			)
			assert.NoError(t, err)
		})
	}
}

func TestNewConfigsFromDir(t *testing.T) {
	_, err := NewFromDir(
		exampleDir,
		WithJustValidation(true),
	)
	assert.NoError(t, err)
}

func TestSchedule_SetDefaults_Schedule(t *testing.T) {
	conf := &Config{}
	defaults(conf)
	assert.Equal(t, TypeSchedulerStrategyTick, conf.Schedule.Strategy)
}

func TestSchedule_Validate(t *testing.T) {
	h := time.Hour
	c := "* * * * *"

	testCases := []struct {
		name     string
		schedule Schedule
		err      error
	}{
		{
			name:     "must_not_miss_interval_and_cron",
			schedule: Schedule{},
			err:      fmt.Errorf("must set schedule.interval or schedule.cron"),
		},
		{
			name: "must_not_have_interval_and_cron",
			schedule: Schedule{
				Interval: &h,
				Cron:     &c,
			},
			err: fmt.Errorf("must set either schedule.interval or schedule.cro"),
		},
		{
			name: "must_have_valid_strategy",
			schedule: Schedule{
				Cron: &c,
			},
			err: fmt.Errorf("unknown strategy: \"\""),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.schedule.Validate()
			if tc.err == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, tc.err, err)
			}
		})
	}
}
