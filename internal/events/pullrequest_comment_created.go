package events

import (
	"github.com/manomartins/bitbird/internal/interfaces"
	"strconv"
)

type PullRequestCommentCreated struct {
	notifier        interfaces.Notifier
	messagesStorage interfaces.PullRequestMessagesInterface
}

func NewPullRequestCommentCreated(notifier interfaces.Notifier, messagesStorage interfaces.PullRequestMessagesInterface) *PullRequestCommentCreated {
	return &PullRequestCommentCreated{
		notifier:        notifier,
		messagesStorage: messagesStorage,
	}
}

func (p *PullRequestCommentCreated) Execute(event PullRequestEvent) error {
	pr, err := p.messagesStorage.GetById(strconv.Itoa(event.PullRequest.ID))
	if err != nil {
		return err
	}

	err = p.notifier.SendCommentNotification(pr.ChannelID, pr.MessageID, event.Comment.Content.Raw)
	if err != nil {
		return err
	}

	return nil

}
