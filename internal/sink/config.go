package sink

import (
	"github.com/turbolytics/latte/internal/collector/template"
)

type Type string

const (
	TypeConsole Type = "console"
	TypeHTTP    Type = "http"
	TypeKafka   Type = "kafka"
	TypeFile    Type = "file"
	TypeS3      Type = "s3"
)

type Config struct {
	Type   Type
	Config map[string]any
}

func ApplyTemplates(c *Config) error {
	// enabling templating across a couple of fixed, known configuration fields
	fields := []string{
		"uri",
	}
	for _, field := range fields {
		if _, hasField := c.Config[field]; hasField {
			bs, err := template.Parse([]byte(c.Config[field].(string)))
			if err != nil {
				return err
			}
			c.Config[field] = string(bs)
		}
	}

	return nil
}
