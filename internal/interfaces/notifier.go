package interfaces

type Notifier interface {
	SendNotification(channelID string, message string) (string, error)
	AddApprovalEmoji(channelID string, messageId string) error
	AddChangeRequestEmoji(channelID string, messageId string) error
	SendCommentNotification(channelID string, messageId string, comment string) error
	RemoveEmoji(channelID string, messageId string) error
}
