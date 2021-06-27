package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v2/bot"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/sanity-io/litter"
	"github.com/sirupsen/logrus"
	"github.com/timkelleher/ferry-bot/pkg/cache"
	"github.com/timkelleher/ferry-bot/pkg/wsdot"
)

type Bot struct {
	Ctx    *bot.Context
	Logger *logrus.Logger
}

func (b *Bot) logEvent(e *gateway.MessageCreateEvent, cmd string) {
	b.Logger.WithFields(logrus.Fields{
		"command":  cmd,
		"author":   e.Author.Username,
		"authorID": e.Author.ID,
		"channel":  e.ChannelID,
		"guild":    e.GuildID,
	}).Info("Message acknowledged")
}

func (b *Bot) Ping(e *gateway.MessageCreateEvent) (string, error) {
	b.logEvent(e, "Ping")
	fmt.Println(e.GuildID)
	return "Pong!", nil
}

func (b *Bot) Terminals(e *gateway.MessageCreateEvent) (string, error) {
	b.logEvent(e, "Terminals")

	endpoint := wsdot.Endpoints[wsdot.TerminalsID]
	payload, err := cache.Get(wsdot.TerminalsID, 24*time.Hour, b.Logger, endpoint.Payload, wsdot.EndpointArguments{
		Date: time.Now(),
	})
	if err != nil {
		b.Logger.WithFields(logrus.Fields{
			"endpoint": wsdot.TerminalsID,
		}).Fatalf("Cannot fetch payload")
	}
	terms, err := endpoint.Unmarshal(payload)
	if err != nil {
		b.Logger.WithFields(logrus.Fields{
			"endpoint": wsdot.TerminalsID,
		}).Fatalf("Cannot parse payload")
	}

	terminals := make([]string, 0)
	terminalsMap := terms.(map[string]wsdot.Terminal)
	for _, terminal := range terminalsMap {
		terminals = append(terminals, terminal.Description)
	}

	return strings.Join(terminals, " | "), nil
}

func (b *Bot) Schedule(e *gateway.MessageCreateEvent, args bot.RawArguments) (string, error) {
	litter.Dump(args)
	return "I dunno", nil
}
