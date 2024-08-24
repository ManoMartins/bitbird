package interfaces

import (
	"context"
	"time"
)

type EmbedField struct {
	Name   string
	Value  string
	Inline bool
}

type EmbedData struct {
	Title     string
	CreatedAt time.Time
	Message   string
	Author    string
	AuthorURL string
	Fields    []*EmbedField
	Content   string
}

type Notifier interface {
	SendNotification(channelID string, message string) (string, error)
	AddApprovalEmoji(channelID string, messageId string) error
	AddChangeRequestEmoji(channelID string, messageId string) error
	SendCommentNotification(channelID string, messageId string, comment string) error
	RemoveEmoji(channelID string, messageId string) error
	SendNotificationEmbed(ctx context.Context, channelID string, embed EmbedData) (string, error)
	GetUserAvatarURL(ctx context.Context, userID string) (string, error)
}
