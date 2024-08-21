package configs

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
)

var dg *discordgo.Session

func ConnectDiscord() {
	token := os.Getenv("DISCORD_BOT_TOKEN")

	dgConnect, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	err = dgConnect.Open()
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}

	log.Println("Connected to Discord...")
	dg = dgConnect
}

func GetDiscord() *discordgo.Session {
	return dg
}

func DisconnectDiscord() {
	if err := dg.Close(); err != nil {
		log.Fatalf("Failed to disconnect Discord: %v", err)
	}

	log.Println("Disconnected from Discord")
}
