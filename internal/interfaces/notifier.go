package interfaces

type Notifier interface {
	SendNotification(message string) (string, error)
	AddApprovalEmoji(messageId string) error
	AddChangeRequestEmoji(messageId string) error
	SendCommentNotification(messageId string, comment string) error
}
