package initializer

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

var exampleDir string

func init() {
	currDir, _ := os.Getwd()
	exampleDir = path.Join(currDir, "..", "..", "..", "dev", "examples")
}

func TestNewConfigFromFile(t *testing.T) {
	testCases := []struct {
		fileName string
	}{
		{"mongo.http.yaml"},
		{"postgres.fileaudit.yaml"},
		{"postgres.http.yaml"},
		{"postgres.kafka.yaml"},
		{"postgres.s3.yaml"},
		{"postgres.stdout.yaml"},
		{"prometheus.fileaudit.yaml"},
	}
	for _, tc := range testCases {
		t.Run(tc.fileName, func(t *testing.T) {
			fPath := path.Join(exampleDir, tc.fileName)
			_, err := NewCollectorFromFile(
				fPath,
				WithJustValidation(true),
			)
			assert.NoError(t, err)
		})
	}
}

func TestNewConfigsFromGlob(t *testing.T) {
	glob := path.Join(exampleDir, "*.yaml")
	cs, err := NewCollectorsFromGlob(
		glob,
		WithJustValidation(true),
	)
	assert.NoError(t, err)
	assert.Equal(t, 7, len(cs))
}
