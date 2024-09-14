package events

import (
	"context"
	"fmt"
	"github.com/manomartins/bitbird/internal/app/interfaces"
	"github.com/manomartins/bitbird/internal/app/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"os"
	"strconv"
	"strings"
	"time"
)

var DiscordUsers = map[string]string{
	"alexandre_valim":                "1258102843798323404",
	"david_almeida_santos":           "550057532631154719",
	"gabriel_alves_de_lima":          "310979218320261120",
	"jean_paes_rabello":              "260392808723120138",
	"joao_victor_pereira_dos_santos": "594647301411176451",
	"liziane_tamm":                   "714928138777526313",
	"manoel_martins":                 "667184274428002345",
	"matheus_de_paula_cordeiro":      "1055305522720555069",
	"marcio_d_carvalho":              "642135148334415882",
	"tassyo_monteiro":                "905438526055800874",
	"william_rodrigues":              "960526439927672902",
	"ana_alice_honorio":              "821798354152456313",
	"luan_s_calais":                  "610988261665538059",
	"samantha_vale":                  "819642515925237790",
	"luiz_amorim":                    "443866985328017408",

	//"henrique_siqueira_cheim": "<@>",
	//"guilherme_borba":         "<@>",
	//"rafael_costa":            "<@>",
	//"islanilton_rodrigues":    "<@>",
}

const name = "github.com/manomartins/bitbird"

var (
	tracer = otel.Tracer(name)
	meter  = otel.Meter(name)
	//logger  = otelslog.NewLogger(name)
	prCounter metric.Int64Counter
)

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

func init() {
	var err error
	prCounter, err = meter.Int64Counter(
		"pull_request_created",
		metric.WithDescription("The number of pull request created"),
		metric.WithUnit("{call}"),
	)
	if err != nil {
		panic(err)
	}
}

func NewPullRequestCreated(notifier interfaces.Notifier, messagesStorage interfaces.PullRequestMessagesInterface) *PullRequestCreated {
	return &PullRequestCreated{
		notifier:        notifier,
		messagesStorage: messagesStorage,
	}
}

func (p *PullRequestCreated) Execute(ctx context.Context, event PullRequestEvent) error {
	opt := metric.WithAttributes(
		attribute.Key("codebase").String(event.Repository.Name),
	)
	prCounter.Add(ctx, 1, opt)

	var prReviewersName []string
	for _, reviewer := range event.PullRequest.Reviewers {
		prReviewersName = append(prReviewersName, reviewer.DisplayName)
	}

	ctx, spanSendMessage := tracer.Start(ctx, "pull_request_created.send_message")
	defer spanSendMessage.End()

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

	avatarURL, err := p.GetUserAvatarURL(ctx, event.Actor.DisplayName)
	if err != nil {
		return err
	}

	users := p.getContent(append([]string{event.Actor.DisplayName}, prReviewersName...))

	channelID := os.Getenv("DISCORD_CHANNEL_ID_FOR_PR")
	messageID, err := p.notifier.SendNotificationEmbed(ctx, channelID, interfaces.EmbedData{
		Title:     "🚀 Novo Pull Request",
		CreatedAt: time.Now(),
		Message:   message,
		Author:    event.Actor.DisplayName,
		AuthorURL: avatarURL,
		Content:   users,
		Fields: []*interfaces.EmbedField{
			{
				Name:   " ",
				Value:  fmt.Sprintf("[**🔗 Clique aqui para ver o PR**](%s)", event.PullRequest.Links.HTML.Href),
				Inline: false,
			},
		},
	})

	if err != nil {
		return err
	}

	err = p.messagesStorage.Create(ctx, strconv.Itoa(event.PullRequest.ID), channelID, messageID)
	if err != nil {
		return err
	}

	return nil
}

func (p *PullRequestCreated) GetUserAvatarURL(ctx context.Context, input string) (string, error) {
	userID, ok := DiscordUsers[utils.ToSnakeCase(input)]

	if ok {
		return p.notifier.GetUserAvatarURL(ctx, userID)
	}

	return "", nil
}

func (p *PullRequestCreated) formatMessage(pr FormatMessageData) string {
	message := fmt.Sprintf("Título: `%s`\n", pr.Title)
	message += fmt.Sprintf("Destino: `%s`\n", pr.Destination)
	message += fmt.Sprintf("Repositorio: `%s`\n", pr.RepoName)

	if len(pr.Reviewers) > 0 {
		message += "\n**📝 Revisores:**\n"
		for _, reviewer := range pr.Reviewers {
			message += fmt.Sprintf("- %s\n", reviewer)
		}
	} else {
		message += "\n*Nenhum revisor atribuído.*\n"
	}

	return message
}

func (p *PullRequestCreated) getContent(users []string) string {
	var usersMention []string
	for _, user := range users {
		userMention, ok := DiscordUsers[utils.ToSnakeCase(user)]

		if ok {
			userMention = fmt.Sprintf("<@%s>", userMention)
			usersMention = append(usersMention, userMention)
		}
	}

	return fmt.Sprintf("📰 Um novo Pull Request foi aberto ||%s||", strings.Join(usersMention, " "))
}
