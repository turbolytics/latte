package sc

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

var configDir string

func init() {
	currDir, _ := os.Getwd()
	configDir = path.Join(currDir, "..", "..", "..", "dev", "scconfigs")
}

func TestNewConfigFromFile(t *testing.T) {
	testCases := []struct {
		fileName string
	}{
		{"sc.memory.yaml"},
	}

	for _, tc := range testCases {
		t.Run(tc.fileName, func(t *testing.T) {
			fPath := path.Join(configDir, tc.fileName)
			_, err := NewFromFile(
				fPath,
			)
			assert.NoError(t, err)
		})
	}
}
