package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/nullspace-lab/discord-bot/config"
	"github.com/nullspace-lab/discord-bot/internal/bot"
)

func main() {
	cfg := config.Load()

	b := bot.New(cfg)
	b.Start()
	defer b.Stop()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
}
