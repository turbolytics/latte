package source

type Type string

const (
	TypeMongoDB     Type = "mongodb"
	TypePartitionS3 Type = "partition.s3"
	TypePostgres    Type = "postgres"
	TypePrometheus  Type = "prometheus"
)

type Config struct {
	Config map[string]any
	Type   Type
}
