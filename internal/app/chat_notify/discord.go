package chat_notify

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/manomartins/bitbird/configs"
	"github.com/manomartins/bitbird/internal/app/interfaces"
	"os"
	"time"
)

type DiscordNotifier struct {
	dg        *discordgo.Session
	channelID string
}

func NewDiscordNotifier() *DiscordNotifier {
	channelID := os.Getenv("DISCORD_CHANNEL_ID")

	dg := configs.GetDiscord()

	return &DiscordNotifier{
		dg:        dg,
		channelID: channelID,
	}
}

func (d *DiscordNotifier) SendNotification(channelID string, message string) (string, error) {
	send, err := d.dg.ChannelMessageSend(channelID, message)

	if err != nil {
		return "", err
	}

	return send.ID, nil
}

func (d *DiscordNotifier) SendCommentNotification(channelID string, messageId string, comment string) error {
	message, err := d.dg.ChannelMessage(channelID, messageId)

	if err != nil {
		return err
	}

	var thread *discordgo.Channel

	if !message.Thread.IsThread() {
		thread, err = d.dg.MessageThreadStart(channelID, messageId, "Novo Topic", 60)
	} else {
		thread = message.Thread
	}

	if err != nil {
		return err
	}

	_, err = d.dg.ChannelMessageSend(thread.ID, comment)
	if err != nil {
		return err
	}

	return nil
}

func (d *DiscordNotifier) AddApprovalEmoji(channelID string, messageID string) error {
	err := d.dg.MessageReactionAdd(channelID, messageID, "✅")

	if err != nil {
		return err
	}

	return nil
}

func (d *DiscordNotifier) AddChangeRequestEmoji(channelID string, messageID string) error {
	err := d.dg.MessageReactionAdd(channelID, messageID, "🔄")

	if err != nil {
		return err
	}

	return nil
}

func (d *DiscordNotifier) RemoveEmoji(channelID string, messageID string) error {
	err := d.dg.MessageReactionsRemoveAll(channelID, messageID)

	if err != nil {
		return err
	}

	return nil
}

func (d *DiscordNotifier) SendNotificationEmbed(ctx context.Context, channelID string, embed interfaces.EmbedData) (string, error) {
	var fields []*discordgo.MessageEmbedField

	for i := range embed.Fields {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   embed.Fields[i].Name,
			Value:  embed.Fields[i].Value,
			Inline: embed.Fields[i].Inline,
		})
	}

	message, err := d.dg.ChannelMessageSendComplex(
		channelID,
		&discordgo.MessageSend{
			Content: embed.Content,
			Embed: &discordgo.MessageEmbed{
				Title:       embed.Title,
				Description: embed.Message,
				Timestamp:   embed.CreatedAt.Format(time.RFC3339),
				Color:       0x8e2cf0,
				Fields:      fields,
				Footer: &discordgo.MessageEmbedFooter{
					Text:    embed.Author,
					IconURL: embed.AuthorURL,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	return message.ID, nil
}

func (d *DiscordNotifier) GetUserAvatarURL(ctx context.Context, userID string) (string, error) {
	user, err := d.dg.User(userID)

	if err != nil {
		return "", err
	}

	avatarURL := user.AvatarURL("")

	return avatarURL, nil
}

func (d *DiscordNotifier) SendDirectMessage(ctx context.Context, userID string, message string) error {
	channel, err := d.dg.UserChannelCreate(userID)
	if err != nil {
		return err
	}

	_, err = d.dg.ChannelMessageSend(channel.ID, message)
	if err != nil {
		return err
	}

	return nil
}
