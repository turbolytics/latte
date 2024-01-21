package mongodb

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func TestParseAgg_MultiStage(t *testing.T) {
	agg := `
      [
        { 
          "$group": { "_id": "$account", "value": { "$count": {} } } 
        },
        { 
          "$addFields": { 
            "account": {  "$toString": "$_id" } 
          } 
        }, 
        { 
          "$project": { "_id": 0 }
        }
      ]
`
	p, err := ParseAgg(agg)
	assert.NoError(t, err)
	assert.Equal(t, mongo.Pipeline{
		primitive.D{
			primitive.E{
				Key: "$group", Value: primitive.D{
					primitive.E{
						Key: "_id", Value: "$account",
					},
					primitive.E{
						Key: "value", Value: primitive.D{
							primitive.E{
								Key: "$count", Value: primitive.D{},
							},
						},
					},
				},
			},
		},
		primitive.D{
			primitive.E{
				Key: "$addFields", Value: primitive.D{
					primitive.E{
						Key: "account", Value: primitive.D{
							primitive.E{
								Key: "$toString", Value: "$_id"},
						},
					},
				},
			},
		},
		primitive.D{
			primitive.E{
				Key: "$project", Value: primitive.D{
					primitive.E{
						Key: "_id", Value: int32(0),
					},
				},
			},
		},
	}, p)
}
