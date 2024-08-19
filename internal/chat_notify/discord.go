package chat_notify

import (
	"github.com/bwmarrin/discordgo"
	"github.com/manomartins/bitbird/configs"
	"os"
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
