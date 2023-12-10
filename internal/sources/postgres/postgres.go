package postgres

import "github.com/mitchellh/mapstructure"

type config struct {
	URI string
}

type Postgres struct{}

func NewFromGenericConfig(m map[string]any) (*Postgres, error) {
	var conf config
	if err := mapstructure.Decode(m, &conf); err != nil {
		return nil, err
	}

	return &Postgres{}, nil
}
