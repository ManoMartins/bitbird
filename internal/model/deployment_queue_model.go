package model

type DeploymentQueueModel struct {
	ID         string `json:"id" bson:"_id,omitempty"`
	CardKey    string `json:"card_key" bson:"card_key"`
	ChannelID  string `json:"channel_id" bson:"channel_id"`
	MessageID  string `json:"message_id" bson:"message_id"`
	CommitHash string `json:"commit_hash" bson:"commit_hash"`
}
