package events

import (
	"context"
	"fmt"
	"github.com/manomartins/bitbird/internal/app/interfaces"
	"github.com/manomartins/bitbird/internal/app/utils"
	"slices"
	"strconv"
)

var acceptDM = []string{"manoel_martins", "jean_paes_rabello", "liziane_tamm", "tassyo_monteiro"}

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

	if !slices.Contains(acceptDM, utils.ToSnakeCase(event.PullRequest.Actor.DisplayName)) {
		return nil
	}

	actorID, ok := DiscordUsers[utils.ToSnakeCase(event.PullRequest.Actor.DisplayName)]
	if !ok {
		return nil
	}

	message := fmt.Sprintf(
		"%s o pull request **%s** foi aprovado! [**Clique aqui para ver o PR**](%s). 🎉",
		event.PullRequest.Actor.DisplayName,
		event.PullRequest.Title,
		event.PullRequest.Links.HTML.Href,
	)

	err = p.notifier.SendDirectMessage(ctx, actorID, message)
	if err != nil {
		return err
	}

	return nil
}
