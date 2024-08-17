package events

import (
	"fmt"
	"github.com/manomartins/bitbird/internal/interfaces"
	"strconv"
)

var discordUsers = map[string]string{
	"manoel_martins":  "<@667184274428002345>",
	"tassyo_monteiro": "<@905438526055800874>",
}

type FormatMessageData struct {
	ID          int
	Title       string
	State       string
	Author      string
	Destination string
	RepoName    string
	Reviewers   []string
	Link        string
}

type PullRequestCreated struct {
	notifier        interfaces.Notifier
	messagesStorage interfaces.PullRequestMessagesInterface
}

func NewPullRequestCreated(notifier interfaces.Notifier, messagesStorage interfaces.PullRequestMessagesInterface) *PullRequestCreated {
	return &PullRequestCreated{
		notifier:        notifier,
		messagesStorage: messagesStorage,
	}
}

func (p *PullRequestCreated) Execute(event PullRequestEvent) error {
	var prReviewersName []string
	for _, review := range event.PullRequest.Reviewers {
		prReviewersName = append(prReviewersName, review.User.DisplayName)
	}

	message := p.formatMessage(
		FormatMessageData{
			ID:          event.PullRequest.ID,
			Title:       event.PullRequest.Title,
			State:       event.PullRequest.State,
			Author:      event.Actor.DisplayName,
			Destination: event.PullRequest.Source.Branch.Name,
			Reviewers:   prReviewersName,
			RepoName:    event.Repository.Name,
			Link:        event.PullRequest.Links.HTML.Href,
		})

	messageID, err := p.notifier.SendNotification(message)
	if err != nil {
		return err
	}

	err = p.messagesStorage.UpdatePullRequestMessage(strconv.Itoa(event.PullRequest.ID), messageID)
	if err != nil {
		return err
	}

	return nil
}

func (p *PullRequestCreated) formatMessage(pr FormatMessageData) string {
	authorMention, ok := discordUsers[pr.Author]
	if !ok {
		authorMention = pr.Author // Caso não haja mapeamento, use o nome do autor
	}

	message := "**🚀 *Detalhes do Pull Request:* **\n\n"
	message += fmt.Sprintf("**Título:** `%s`\n", pr.Title)
	message += fmt.Sprintf("**Status:** `%s`\n", pr.State)
	message += fmt.Sprintf("**Autor:** %s\n", authorMention)
	message += fmt.Sprintf("**Destino:** `%s`\n", pr.Destination)
	message += fmt.Sprintf("**🌟 Repositorio:** `%s`\n", pr.RepoName)
	message += fmt.Sprintf("**Link do PR:** [Abrir PR](%s)\n", pr.Link) // Link do PR

	if len(pr.Reviewers) > 0 {
		message += "\n**📝 Revisores:**\n"
		for _, review := range pr.Reviewers {
			reviewerMention, ok := discordUsers[review]
			if !ok {
				reviewerMention = review // Caso não haja mapeamento
			}
			message += fmt.Sprintf("- %s\n", reviewerMention)
		}
	} else {
		message += "\n*Nenhum revisor atribuído.*\n"
	}

	return message
}
