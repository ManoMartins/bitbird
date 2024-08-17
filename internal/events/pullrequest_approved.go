package events

import (
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

func (p *PullRequestApproved) Execute(event PullRequestEvent) error {
	pr, err := p.messagesStorage.GetPullRequestMessage(strconv.Itoa(event.PullRequest.ID))

	if err != nil {
		return err
	}

	err = p.notifier.AddApprovalEmoji(pr.MessageID)
	if err != nil {
		return err
	}

	return nil
}
