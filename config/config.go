package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Token               string
	GuildID             string
	RulesChannelID      string
	NewMemberRoleID     string
	OficialMemberRoleID string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		Token:               os.Getenv("DISCORD_TOKEN"),
		GuildID:             os.Getenv("DISCORD_GUILD_ID"),
		RulesChannelID:      os.Getenv("DISCORD_RULES_CHANNEL_ID"),
		NewMemberRoleID:     os.Getenv("DISCORD_NEW_MEMBER_ROLE_ID"),
		OficialMemberRoleID: os.Getenv("DISCORD_OFICIAL_MEMBER_ROLE_ID"),
	}
}
