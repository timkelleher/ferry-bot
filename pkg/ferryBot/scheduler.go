package ferryBot

import (
	"context"
	"time"

	"github.com/diamondburned/arikawa/v2/state"
	"github.com/sirupsen/logrus"
)

type SchedulerConfig struct {
	Name     string
	Logger   *logrus.Logger
	Duration time.Duration
	State    *state.State
}

func NewScheduler(ctx context.Context, conf SchedulerConfig, fn func(*state.State) error) func() {
	return func() {
		ticker := time.NewTicker(conf.Duration)
		for {
			select {
			case <-ctx.Done():
				conf.Logger.WithFields(logrus.Fields{
					"name": conf.Name,
				}).Info("Scheduler has stopped")

				return
			case <-ticker.C:
				conf.Logger.WithFields(logrus.Fields{
					"name": conf.Name,
				}).Info("Scheduler is running")

				err := fn(conf.State)
				if err != nil {
					conf.Logger.WithFields(logrus.Fields{
						"name":  conf.Name,
						"error": err.Error(),
					}).Error("Scheduler encountered an error")
				}
			}
		}
	}
}
