package events

import (
	"context"
	"github.com/manomartins/bitbird/internal/interfaces"
	"strconv"
)

type PullRequestApproved struct {
	notifier        interfaces.Notifier
	messagesStorage interfaces.PullRequestMessagesInterface
}

func NewPullRequestApproved(notifier interfaces.Notifier, messagesStorage interfaces.PullRequestMessagesInterface) *PullRequestApproved {
	return &PullRequestApproved{
		notifier:        notifier,
		messagesStorage: messagesStorage,
	}
}

func (p *PullRequestApproved) Execute(ctx context.Context, event PullRequestEvent) error {
	pr, err := p.messagesStorage.GetById(strconv.Itoa(event.PullRequest.ID))

	if err != nil {
		return err
	}

	err = p.notifier.AddApprovalEmoji(pr.ChannelID, pr.MessageID)
	if err != nil {
		return err
	}

	return nil
}
