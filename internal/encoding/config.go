package encoding

import (
	"bytes"
	"github.com/turbolytics/latte/internal/encoding/json"
	"github.com/turbolytics/latte/internal/encoding/parquet"
)

type Type string

const (
	TypeParquet Type = "parquet"
	TypeJSON    Type = "json"
)

type Config struct {
	Type   Type
	Config map[string]any
}

type Encoder interface {
	Init(*bytes.Buffer) error
	Write(any) error
	Flush() error
	Close() error
}

func NewEncoder(c Config) (Encoder, error) {
	var err error
	var e Encoder

	switch c.Type {
	case TypeParquet:
		e, err = parquet.NewFromGenericConfig(c.Config)
	default:
		e, err = json.NewFromGenericConfig(c.Config)
	}
	return e, err
}
