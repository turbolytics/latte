package parquet

import (
	"bytes"
	"github.com/mitchellh/mapstructure"
)

type Field struct {
	Name string
	Type string
	From string
}

type config struct {
	Compression string
	Schema      []Field
}

type Parquet struct {
	config config
}

func (p *Parquet) Init(buf *bytes.Buffer) error {
	return nil
}

func (p *Parquet) Write(d any) error {
	return nil
}

func (p *Parquet) Close() error {
	return nil
}

func (p *Parquet) Flush() error {
	return nil
}

func NewFromGenericConfig(m map[string]any) (*Parquet, error) {
	var conf config
	if err := mapstructure.Decode(m, &conf); err != nil {
		return nil, err
	}

	p := &Parquet{
		config: conf,
	}

	return p, nil
}
