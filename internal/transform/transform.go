package transform

import (
	"bytes"
	"github.com/mitchellh/mapstructure"
	"github.com/turbolytics/latte/internal/record"
	"text/template"
)

type Transformer interface {
	Transform(tContext any) ([]byte, error)
}

type Template struct {
	Template string
}

func (t *Template) Transform(tContext any) ([]byte, error) {
	textTempl, err := template.New("transform").Parse(string(t.Template))
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer

	if err := textTempl.Execute(&out, tContext); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func NewTemplateFromGenericConfig(m map[string]any) (*Template, error) {
	var t Template
	if err := mapstructure.Decode(m, &t); err != nil {
		return nil, err
	}

	return &t, nil
}

type Noop struct{}

func (n Noop) Transform(record.Result) error {
	return nil
}
