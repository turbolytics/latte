package mongodb

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/turbolytics/latte/internal/metric"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func TestIntegration_Mongo_Source(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	mongodbContainer, err := mongodb.RunContainer(
		ctx,
		testcontainers.WithImage("mongo:6"),
	)

	// Clean up the container
	defer func() {
		if err := mongodbContainer.Terminate(ctx); err != nil {
			assert.NoError(t, err)
		}
	}()

	endpoint, err := mongodbContainer.ConnectionString(ctx)
	if err != nil {
		assert.NoError(t, err)
	}

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(endpoint))
	if err != nil {
		assert.NoError(t, err)
	}

	type user struct {
		Account    string
		SignupTime time.Time
	}

	coll := mongoClient.Database("test").Collection("users")
	docs := []interface{}{
		user{Account: "amazon", SignupTime: time.Now().UTC()},
		user{Account: "amazon", SignupTime: time.Now().UTC()},
		user{Account: "google", SignupTime: time.Now().UTC()},
	}
	result, err := coll.InsertMany(context.TODO(), docs)
	fmt.Printf("Documents inserted: %v\n", len(result.InsertedIDs))

	mSource, err := NewFromGenericConfig(
		ctx,
		map[string]any{
			"uri":        endpoint,
			"database":   "test",
			"collection": "users",
			"agg": `
  [
	{ 
	  "$group": { "_id": "$account", "value": { "$count": {} } } 
	},
	{ "$sort" : { "_id" : 1 } },
	{ 
	  "$addFields": { 
		"account": {  "$toString": "$_id" } 
	  } 
	}, 
	{ 
	  "$project": { "_id": 0 }
	}
  ]
`,
		},
		false,
	)

	assert.NoError(t, err)
	rs, err := mSource.Source(ctx)
	assert.NoError(t, err)

	ms := rs.(*metric.Metrics)

	for _, m := range ms.Metrics {
		m.Timestamp = time.Time{}
		m.UUID = ""
	}

	assert.Equal(t, []*metric.Metric{
		{
			Value: 2,
			Tags: map[string]string{
				"account": "amazon",
			},
		},
		{
			Value: 1,
			Tags: map[string]string{
				"account": "google",
			},
		},
	}, ms.Metrics)

}
