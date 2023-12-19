package mongodb

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Test struct{}

func ParseAgg(agg string) (mongo.Pipeline, error) {
	var q interface{}
	bs := []byte(agg)
	if err := bson.UnmarshalExtJSON(bs, true, &q); err != nil {
		return nil, err
	}

	var p mongo.Pipeline
	for _, stage := range q.(bson.A) {
		p = append(p, stage.(bson.D))
	}

	return p, nil
}
