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

/*
func (c *Collector) Init(l *zap.Logger, validate bool) error {
	var err error
	var s Sourcer
	switch c.Type {
			case TypePartitionS3:
				s, err = s3.NewFromGenericConfig(
					c.Collector,
				)
			}
		case TypePostgres:
			s, err = postgres.NewFromGenericConfig(
				c.Collector,
				validate,
			)
	case TypeMongoDB:
		s, err = mongodb.NewFromGenericConfig(
			context.TODO(),
			c.Collector,
			validate,
		)
			case TypePrometheus:
				s, err = prometheus.NewFromGenericConfig(
					c.Collector,
					prometheus.WithLogger(l),
				)
	default:
		return fmt.Errorf("source type: %q unknown", c.Type)
	}

	if err != nil {
		return err
	}
	c.Sourcer = s
	return nil
}
*/
