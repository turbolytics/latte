package collector

import (
	"os"
	"path"
)

var exampleDir string

func init() {
	currDir, _ := os.Getwd()
	exampleDir = path.Join(currDir, "..", "..", "dev", "examples")
}

/*
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
		{"s3.http.clickhouse.yaml"},
	}
	for _, tc := range testCases {
		t.Run(tc.fileName, func(t *testing.T) {
			fPath := path.Join(exampleDir, tc.fileName)
			_, err := invoker.NewFromFile(
				fPath,
				WithJustValidation(true),
			)
			assert.NoError(t, err)
		})
	}
}

func TestNewConfigsFromGlob(t *testing.T) {
	glob := path.Join(exampleDir, "*.yaml")
	cs, err := invoker.NewFromGlob(
		glob,
		WithJustValidation(true),
	)
	assert.NoError(t, err)
	assert.Equal(t, 7, len(cs))
}
*/
