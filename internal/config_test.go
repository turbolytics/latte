package internal

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

var exampleDir string

func init() {
	currDir, _ := os.Getwd()
	exampleDir = path.Join(currDir, "..", "dev", "examples")
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
			_, err := NewConfigFromFile(
				fPath,
				WithJustConfigValidation(true),
			)
			assert.NoError(t, err)
		})
	}
}

func TestNewConfigsFromDir(t *testing.T) {
	_, err := NewConfigsFromDir(
		exampleDir,
		WithJustConfigValidation(true),
	)
	assert.NoError(t, err)
}

func TestSchedule_SetDefaults_Schedule(t *testing.T) {
	conf := &Config{}
	defaults(conf)
	assert.Equal(t, TypeSchedulerStrategyTick, conf.Schedule.Strategy)
}
