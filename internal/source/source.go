package source

type Type string

const (
	TypeMetricMongoDB  Type = "metric.mongodb"
	TypeMetricPostgres Type = "metric.postgres"
	// TypePartitionS3    Type = "partition.s3"
	TypePrometheus Type = "prometheus"
)

type Config struct {
	Config map[string]any
	Type   Type
}
