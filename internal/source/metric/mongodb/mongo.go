package mongodb

import (
	"context"
	"github.com/mitchellh/mapstructure"
	"github.com/turbolytics/latte/internal/metric"
	"github.com/turbolytics/latte/internal/record"
	"github.com/turbolytics/latte/internal/source"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type config struct {
	URI        string
	Agg        string
	Database   string
	Collection string
}

type Mongo struct {
	config config
	client *mongo.Client
}

func (m *Mongo) WindowDuration() *time.Duration {
	return nil
}

func (m Mongo) Type() source.Type {
	return source.TypeMetricMongoDB
}

func (m *Mongo) Source(ctx context.Context) (record.Result, error) {
	p, err := ParseAgg(m.config.Agg)
	if err != nil {
		return nil, err
	}

	col := m.client.Database(m.config.Database).Collection(m.config.Collection)
	cursor, err := col.Aggregate(
		context.TODO(),
		p,
	)
	if err != nil {
		return nil, err
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	var rs []map[string]any
	for _, r := range results {
		rs = append(rs, r)
	}

	ms, err := metric.MapsToMetrics(rs)

	metricsResult := metric.NewMetricsResult(ms)
	return metricsResult, err
}

func NewFromGenericConfig(ctx context.Context, m map[string]any, validate bool) (*Mongo, error) {
	var conf config
	if err := mapstructure.Decode(m, &conf); err != nil {
		return nil, err
	}

	var client *mongo.Client
	var err error
	if validate {
		if _, err = ParseAgg(conf.Agg); err != nil {
			return nil, err
		}
	} else {
		client, err = mongo.Connect(ctx, options.Client().ApplyURI(conf.URI))
		if err != nil {
			return nil, err
		}

		if err = client.Ping(ctx, readpref.SecondaryPreferred()); err != nil {
			return nil, err
		}
	}

	return &Mongo{
		config: conf,
		client: client,
	}, nil
}
