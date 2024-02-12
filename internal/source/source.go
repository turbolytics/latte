package source

type Type string

const (
	TypeMongoDB     Type = "mongodb"
	TypePartitionS3 Type = "partition.s3"
	TypePostgres    Type = "postgres"
	TypePrometheus  Type = "prometheus"
)

type Config struct {
	Config map[string]any
	Type   Type
}

func (c Config) Validate() error {
	/*
		vs := map[TypeStrategy]struct{}{
			TypeStrategyTick:                   {},
			TypeStrategyHistoricTumblingWindow: {},
			TypeStrategyIncremental:            {},
		}

		if _, ok := vs[c.Strategy]; !ok {
			return fmt.Errorf("unknown strategy: %q", c.Strategy)
		}
	*/
	return nil
}

func (c *Config) SetDefaults() {
	/*
		if c.Strategy == "" {
			c.Strategy = TypeStrategyTick
		}
	*/
}
