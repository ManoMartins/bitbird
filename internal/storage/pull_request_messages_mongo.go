package storage

import (
	"context"
	"github.com/manomartins/bitbird/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PullRequestMessagesMongo struct {
	collection *mongo.Collection
}

func NewPullRequestMessagesMongo() *PullRequestMessagesMongo {
	mongoClient := GetMongoClient()
	collection := mongoClient.Database("bitbard").Collection("pull_request_message")

	return &PullRequestMessagesMongo{
		collection: collection,
	}
}

func (p *PullRequestMessagesMongo) GetById(prID string) (*model.PullRequestMessageModel, error) {
	var pr model.PullRequestMessageModel

	err := p.collection.FindOne(context.TODO(),
		bson.D{
			{"pr_id", prID},
		}).Decode(&pr)

	if err != nil {
		return nil, err
	}

	return &pr, nil
}

func (p *PullRequestMessagesMongo) FindAll() ([]model.PullRequestMessageModel, error) {
	var prs []model.PullRequestMessageModel

	cursor, err := p.collection.Find(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.TODO(), &prs)
	if err != nil {
		return nil, err
	}

	return prs, nil
}

func (p *PullRequestMessagesMongo) Save() {
	//TODO implement me
	panic("implement me")
}

func (p *PullRequestMessagesMongo) Create(ctx context.Context, prID string, channelID string, messageID string) error {
	_, err := p.collection.InsertOne(
		ctx,
		model.PullRequestMessageModel{
			PrID:      prID,
			ChannelID: channelID,
			MessageID: messageID,
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func (p *PullRequestMessagesMongo) Update(prID string, channelID string, messageID string) error {
	_, err := p.collection.InsertOne(
		context.TODO(),
		model.PullRequestMessageModel{
			PrID:      prID,
			ChannelID: channelID,
			MessageID: messageID,
		},
	)

	if err != nil {
		return err
	}

	return nil
}
