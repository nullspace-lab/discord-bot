package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/nullspace-lab/discord-bot/config"
	"github.com/nullspace-lab/discord-bot/internal/handlers"
)

type Bot struct {
	Session *discordgo.Session
	Config  *config.Config
}

func New(cfg *config.Config) *Bot {
	dg, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		log.Fatalf("failed to create Discord session: %v", err)
	}

	dg.Identify.Intents = discordgo.IntentsGuilds |
		discordgo.IntentsGuildMembers |
		discordgo.IntentsGuildMessageReactions

	return &Bot{
		Session: dg,
		Config:  cfg,
	}
}

func (b *Bot) Start() {
	b.Session.AddHandler(handlers.InteractionHandler)
	b.Session.AddHandler(handlers.ReactionAddHandler(b.Config))

	err := b.Session.Open()
	if err != nil {
		log.Fatalf("failed to open Discord connection: %v", err)
	}

	_, err = b.Session.ApplicationCommandCreate(b.Session.State.User.ID, b.Config.GuildID, &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Verifica se o bot está online",
	})
	if err != nil {
		log.Fatalf("failed to register ping command: %v", err)
	}

	handlers.SetupRules(b.Session, b.Config)
	handlers.StartScheduler(b.Session, b.Config)

	log.Println("bot is online")
}

func (b *Bot) Stop() {
	b.Session.Close()
}
