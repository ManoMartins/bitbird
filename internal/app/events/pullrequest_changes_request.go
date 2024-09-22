package events

import (
	"context"
	"fmt"
	"github.com/manomartins/bitbird/internal/app/interfaces"
	"github.com/manomartins/bitbird/internal/app/utils"
	"slices"
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

func (p *PullRequestChangesRequest) Execute(ctx context.Context, event PullRequestEvent) error {
	pr, err := p.messagesStorage.GetById(strconv.Itoa(event.PullRequest.ID))
	if err != nil {
		return err
	}

	err = p.notifier.AddChangeRequestEmoji(pr.ChannelID, pr.MessageID)
	if err != nil {
		return err
	}

	if !slices.Contains(acceptDM, utils.ToSnakeCase(event.PullRequest.Author.DisplayName)) {
		return nil
	}

	actorID, ok := DiscordUsers[utils.ToSnakeCase(event.PullRequest.Author.DisplayName)]
	if !ok {
		return nil
	}

	message := fmt.Sprintf(
		"%s, foi solicitado que você faça algumas alterações no pull request **%s**. [**Clique aqui para ver os detalhes**](%s). ✏️",
		event.PullRequest.Author.DisplayName,
		event.PullRequest.Title,
		event.PullRequest.Links.HTML.Href,
	)

	err = p.notifier.SendDirectMessage(ctx, actorID, message)
	if err != nil {
		return err
	}

	return nil
}
