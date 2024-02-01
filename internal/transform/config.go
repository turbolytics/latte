package transform

type Type string

const (
	TypeTemplate Type = "template"
)

type Config struct {
	Name        string
	Type        Type
	Template    string
	Config      map[string]any
	Transformer Transformer
}

func (c *Config) Init() error {
	switch c.Type {
	case TypeTemplate:
		trans, err := NewTemplateFromGenericConfig(c.Config)
		if err != nil {
			return err
		}
		c.Transformer = trans
	}

	return nil
}
