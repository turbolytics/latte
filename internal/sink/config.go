package sink

type Config struct {
	Type   Type
	Config map[string]any
}

/*
func (c *Collector) Init(validate bool, logger *zap.Logger) error {
	var sink _type.Sinker
	var err error

	switch c.Type {
	case TypeConsole:
		sink, err = console.NewFromGenericConfig(c.Collector)
	case TypeKafka:
		sink, err = kafka.NewFromGenericConfig(c.Collector)
	case TypeHTTP:
		sink, err = http.NewFromGenericConfig(
			c.Collector,
			http.WithLogger(logger),
		)
	case TypeFile:
		sink, err = file.NewFromGenericConfig(
			c.Collector,
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
