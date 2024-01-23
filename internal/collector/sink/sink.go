package sink

import "github.com/turbolytics/collector/internal/sinks"

type Sink struct {
	Type   sinks.Type
	Sinker sinks.Sinker
	Config map[string]any
}
