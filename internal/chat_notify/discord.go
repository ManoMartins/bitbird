package chat_notify

import (
	"github.com/bwmarrin/discordgo"
	"github.com/manomartins/bitbird/configs"
	"os"
)

type DiscordNotifier struct {
	Dg        *discordgo.Session
	ChannelID string
}

func NewDiscordNotifier() *DiscordNotifier {
	channelID := os.Getenv("DISCORD_CHANNEL_ID")

	dg := configs.GetDiscord()

	return &DiscordNotifier{
		Dg:        dg,
		ChannelID: channelID,
	}
}

func (d *DiscordNotifier) SendNotification(message string) (string, error) {
	send, err := d.Dg.ChannelMessageSend(d.ChannelID, message)

	if err != nil {
		return "", err
	}

	return send.ID, nil
}

func (d *DiscordNotifier) SendCommentNotification(messageId string, comment string) error {
	message, err := d.Dg.ChannelMessage(d.ChannelID, messageId)

	if err != nil {
		return err
	}

	var thread *discordgo.Channel

	if !message.Thread.IsThread() {
		thread, err = d.Dg.MessageThreadStart(d.ChannelID, messageId, "Novo Topic", 60)
	} else {
		thread = message.Thread
	}

	if err != nil {
		return err
	}

	_, err = d.Dg.ChannelMessageSend(thread.ID, comment)
	if err != nil {
		return err
	}

	return nil
}

func (d *DiscordNotifier) AddApprovalEmoji(messageID string) error {
	err := d.Dg.MessageReactionAdd(d.ChannelID, messageID, "✅")

	if err != nil {
		return err
	}

	return nil
}

func (d *DiscordNotifier) AddChangeRequestEmoji(messageID string) error {
	err := d.Dg.MessageReactionAdd(d.ChannelID, messageID, "🔄")

	if err != nil {
		return err
	}

	return nil
}
