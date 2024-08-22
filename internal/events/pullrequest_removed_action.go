package events

import (
	"context"
	"github.com/manomartins/bitbird/internal/interfaces"
	"strconv"
)

type PullRequestRemovedAction struct {
	notifier        interfaces.Notifier
	messagesStorage interfaces.PullRequestMessagesInterface
}

func NewPullRequestRemovedAction(notifier interfaces.Notifier, messagesStorage interfaces.PullRequestMessagesInterface) *PullRequestRemovedAction {
	return &PullRequestRemovedAction{
		notifier:        notifier,
		messagesStorage: messagesStorage,
	}
}

func (p *PullRequestRemovedAction) Execute(ctx context.Context, event PullRequestEvent) error {
	pr, err := p.messagesStorage.GetById(strconv.Itoa(event.PullRequest.ID))

	if err != nil {
		return err
	}

	err = p.notifier.RemoveEmoji(pr.ChannelID, pr.MessageID)
	if err != nil {
		return err
	}

	return nil
}
