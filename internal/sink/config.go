package sink

type Config struct {
	Type   Type
	Config map[string]any
}

/*
func (c *Config) Init(validate bool, logger *zap.Logger) error {
	var sink _type.Sinker
	var err error

	switch c.Type {
	case TypeConsole:
		sink, err = console.NewFromGenericConfig(c.Config)
	case TypeKafka:
		sink, err = kafka.NewFromGenericConfig(c.Config)
	case TypeHTTP:
		sink, err = http.NewFromGenericConfig(
			c.Config,
			http.WithLogger(logger),
		)
	case TypeFile:
		sink, err = file.NewFromGenericConfig(
			c.Config,
			validate,
		)
	}

	if err != nil {
		return err
	}

	c.Sinker = sink
	return nil
}
*/
