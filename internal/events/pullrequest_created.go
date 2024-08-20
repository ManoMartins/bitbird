package events

import (
	"fmt"
	"github.com/manomartins/bitbird/internal/interfaces"
	"github.com/manomartins/bitbird/internal/utils"
	"os"
	"strconv"
)

var DiscordUsers = map[string]string{
	"alexandre_valim":                "<@1258102843798323404>",
	"david_almeida_santos":           "<@550057532631154719>",
	"gabriel_alves_de_lima":          "<@310979218320261120>",
	"jean_paes_rabello":              "<@260392808723120138>",
	"joao_victor_pereira_dos_santos": "<@594647301411176451>",
	"liziane_tamm":                   "<@714928138777526313>",
	"manoel_martins":                 "<@667184274428002345>",
	"matheus_de_paula_cordeiro":      "<@1055305522720555069>",
	"marcio_d_carvalho":              "<@642135148334415882>",
	"tassyo_monteiro":                "<@905438526055800874>",
	"william_rodrigues":              "<@960526439927672902>",
	"ana_alice_honorio":              "<@821798354152456313>",
	"luan_s_calais":                  "<@610988261665538059>",
	"samantha_vale":                  "<@819642515925237790>",
	"luiz_amorim":                    "<@443866985328017408>",

	//"henrique_siqueira_cheim": "<@>",
	//"guilherme_borba":         "<@>",
	//"rafael_costa":            "<@>",
	//"islanilton_rodrigues":    "<@>",
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
	for _, reviewer := range event.PullRequest.Reviewers {
		prReviewersName = append(prReviewersName, reviewer.DisplayName)
	}

	message := p.formatMessage(
		FormatMessageData{
			ID:          event.PullRequest.ID,
			Title:       event.PullRequest.Title,
			State:       event.PullRequest.State,
			Author:      event.Actor.DisplayName,
			Destination: event.PullRequest.Destination.Branch.Name,
			Reviewers:   prReviewersName,
			RepoName:    event.Repository.Name,
			Link:        event.PullRequest.Links.HTML.Href,
		})

	channelID := os.Getenv("DISCORD_CHANNEL_ID_FOR_PR")
	messageID, err := p.notifier.SendNotification(channelID, message)
	if err != nil {
		return err
	}

	err = p.messagesStorage.Create(strconv.Itoa(event.PullRequest.ID), channelID, messageID)
	if err != nil {
		return err
	}

	return nil
}

func (p *PullRequestCreated) formatMessage(pr FormatMessageData) string {
	authorMention, ok := DiscordUsers[utils.ToSnakeCase(pr.Author)]
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
		for _, reviewer := range pr.Reviewers {
			reviewerMention, ok := DiscordUsers[utils.ToSnakeCase(reviewer)]
			if !ok {
				reviewerMention = reviewer // Caso não haja mapeamento
			}
			message += fmt.Sprintf("- %s\n", reviewerMention)
		}
	} else {
		message += "\n*Nenhum revisor atribuído.*\n"
	}

	return message
}
