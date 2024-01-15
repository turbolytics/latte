package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
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

func TestNewConfigsFromGlob(t *testing.T) {
	glob := path.Join(exampleDir, "*.yaml")
	cs, err := NewFromGlob(
		glob,
		WithJustValidation(true),
	)
	assert.NoError(t, err)
	assert.Equal(t, 6, len(cs))
}
