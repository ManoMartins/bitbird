package events

import (
	"github.com/manomartins/bitbird/internal/interfaces"
	"strconv"
)

type PullRequestChangesRequest struct {
	notifier        interfaces.Notifier
	messagesStorage interfaces.PullRequestMessagesInterface
}

func NewPullRequestChangesRequest(notifier interfaces.Notifier, messagesStorage interfaces.PullRequestMessagesInterface) *PullRequestChangesRequest {
	return &PullRequestChangesRequest{
		notifier:        notifier,
		messagesStorage: messagesStorage,
	}
}

func (p *PullRequestChangesRequest) Execute(event PullRequestEvent) error {
	pr, err := p.messagesStorage.GetById(strconv.Itoa(event.PullRequest.ID))
	if err != nil {
		return err
	}

	err = p.notifier.AddChangeRequestEmoji(pr.ChannelID, pr.MessageID)
	if err != nil {
		return err
	}

	return nil
}
