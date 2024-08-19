package model

type PullRequestMessageModel struct {
	ID        string `json:"id" bson:"_id,omitempty"`
	PrID      string `json:"pr_id" bson:"pr_id"`
	ChannelID string `json:"channel_id" bson:"channel_id"`
	MessageID string `json:"message_id" bson:"message_id"`
}
