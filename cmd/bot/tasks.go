package main

import (
	"context"
	"time"

	"github.com/diamondburned/arikawa/v2/state"
	"github.com/sirupsen/logrus"
	"github.com/timkelleher/ferry-bot/pkg/ferryBot"
)

func tasks(ctx context.Context, logger *logrus.Logger, s *state.State) {
	heartbeat := ferryBot.NewScheduler(ctx, ferryBot.SchedulerConfig{
		Name:     "Heartbeat",
		Logger:   logger,
		Duration: 10 * time.Minute,
		State:    s,
	}, func(*state.State) error {
		return nil
	})
	go heartbeat()

	categorySync := ferryBot.NewScheduler(ctx, ferryBot.SchedulerConfig{
		Name:     "CategorySync",
		Logger:   logger,
		Duration: 20 * time.Minute,
		State:    s,
	}, func(s *state.State) error {
		manager := ferryBot.NewChannelManager(s, logger)
		terminals := ferryBot.GetTerminals(logger)
		err := manager.SyncCategories(terminals)
		if err != nil {
			logger.Error("Failed to sync categories")
		}
		return nil
	})
	go categorySync()

	channelSync := ferryBot.NewScheduler(ctx, ferryBot.SchedulerConfig{
		Name:     "ChannelSync",
		Logger:   logger,
		Duration: 15 * time.Second,
		State:    s,
	}, func(s *state.State) error {
		manager := ferryBot.NewChannelManager(s, logger)
		//todo: routes
		err := manager.SyncChannels()
		if err != nil {
			logger.Error("Failed to sync channels")
		}
		return nil
	})
	go channelSync()

	select {}
}
