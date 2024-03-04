package source

import (
	"github.com/turbolytics/latte/internal/collector/template"
)

type Type string

const (
	TypeMetricMongoDB  Type = "metric.mongodb"
	TypeMetricPostgres Type = "metric.postgres"
	TypePrometheus     Type = "metric.prometheus"
	TypeMetricS3       Type = "metric.s3"
)

type Config struct {
	Config map[string]any
	Type   Type
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
