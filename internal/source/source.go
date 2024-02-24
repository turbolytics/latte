package source

type Type string

const (
	TypeMetricMongoDB  Type = "metric.mongodb"
	TypeMetricPostgres Type = "metric.postgres"
	TypePrometheus     Type = "metric.prometheus"
)

type Config struct {
	Config map[string]any
	Type   Type
}
