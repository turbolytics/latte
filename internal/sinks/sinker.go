package sinks

type Type string

const (
	TypeConsole Type = "console"
	TypeHTTP    Type = "http"
)

// Sinker is responsible for sinking
// TODO - Starting with an io.Writer for right now.
type Sinker interface {
	Write([]byte) (int, error)
}
