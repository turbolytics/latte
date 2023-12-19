package mongodb

import (
	"context"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/turbolytics/collector/internal/metrics"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"strconv"
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

func resultsToMetrics(results []bson.M) ([]*metrics.Metric, error) {
	var ms []*metrics.Metric
	for _, r := range results {
		val, ok := r["value"]
		if !ok {
			return nil, fmt.Errorf("each row must contain a %q key", "value")
		}

		m := metrics.New()

		switch v := val.(type) {
		case int:
			m.Value = float64(v)
		case int32:
			m.Value = float64(v)
		case int64:
			m.Value = float64(v)
		case string:
			tv, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, fmt.Errorf("unable to parse string to float: %q", v)
			}
			m.Value = tv
		}
		delete(r, "value")
		for k, v := range r {
			m.Tags[k] = v.(string)
		}
		ms = append(ms, &m)
	}
	return ms, nil
}

func (m *Mongo) Source(ctx context.Context) ([]*metrics.Metric, error) {
	p, err := ParseAgg(m.config.Agg)

	col := m.client.Database(m.config.Database).Collection(m.config.Collection)
	cursor, err := col.Aggregate(context.TODO(), p)
	if err != nil {
		return nil, err
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	ms, err := resultsToMetrics(results)
	return ms, err
}

func NewFromGenericConfig(ctx context.Context, m map[string]any, validate bool) (*Mongo, error) {
	var conf config
	if err := mapstructure.Decode(m, &conf); err != nil {
		return nil, err
	}

	var client *mongo.Client
	var err error
	if validate {
		_, err = ParseAgg(conf.Agg)
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
