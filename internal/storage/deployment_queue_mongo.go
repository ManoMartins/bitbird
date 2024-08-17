package storage

import (
	"context"
	"errors"
	"github.com/manomartins/bitbird/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DeploymentQueueMongo struct {
	collection *mongo.Collection
}

func NewDeploymentQueueMongo() *DeploymentQueueMongo {
	collection := GetMongoClient().Database("bitbard").Collection("deployment_queue")

	return &DeploymentQueueMongo{
		collection: collection,
	}
}

func (d *DeploymentQueueMongo) Create(data model.DeploymentQueueModel) error {
	_, err := d.collection.InsertOne(context.TODO(), data)

	return err
}

func (d *DeploymentQueueMongo) GetByCardKey(key string) (*model.DeploymentQueueModel, error) {
	var result model.DeploymentQueueModel
	err := d.collection.FindOne(context.TODO(), bson.D{{"card_key", key}}).Decode(&result)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &result, err
}
