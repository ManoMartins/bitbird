package model

type PullRequestMessageModel struct {
	ID        string `json:"id" bson:"_id"`
	PrID      string `json:"pr_id" bson:"pr_id"`
	MessageID string `json:"message_id" bson:"message_id"`
}
