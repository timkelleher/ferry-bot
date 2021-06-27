package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/diamondburned/arikawa/v2/bot"
	"github.com/diamondburned/arikawa/v2/state"
	"github.com/sirupsen/logrus"
)

func NewLogger(fileName string) *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("Cannot launch logrus")
	}
	log.Out = file

	return log
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()
	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()

	logger := NewLogger("ferry.log")
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		logger.Fatal("Invalid Discord Token (via env)")
	}

	state, err := state.New("Bot " + token)
	if err != nil {
		logger.WithField("error", err.Error()).Fatal("Failed to connect to Discord")
	}
	//s.SendMessage(858262521227640852, "fffff", nil)

	go tasks(ctx, logger, state)

	logger.Info("Starting Discord Bot")
	bot.Run(token, &Bot{Logger: logger},
		func(ctx *bot.Context) error {
			ctx.HasPrefix = bot.NewPrefix("!")
			return nil
		},
	)
}
