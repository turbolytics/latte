package sink

type Type string

const (
	TypeConsole Type = "console"
	TypeHTTP    Type = "http"
	TypeKafka   Type = "kafka"
	TypeFile    Type = "file"
)

type Config struct {
	Type   Type
	Config map[string]any
}
