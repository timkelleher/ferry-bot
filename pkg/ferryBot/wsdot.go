package ferryBot

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/timkelleher/ferry-bot/pkg/cache"
	"github.com/timkelleher/ferry-bot/pkg/wsdot"
)

func GetTerminals(logger *logrus.Logger) []wsdot.Terminal {
	endpoint := wsdot.Endpoints[wsdot.TerminalsID]
	payload, err := cache.Get(wsdot.TerminalsID, 24*time.Hour, logger, endpoint.Payload, wsdot.EndpointArguments{
		Date: time.Now(),
	})
	if err != nil {
		logger.WithFields(logrus.Fields{
			"endpoint": wsdot.TerminalsID,
		}).Fatalf("Cannot fetch payload")
	}
	terms, err := endpoint.Unmarshal(payload)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"endpoint": wsdot.TerminalsID,
		}).Fatalf("Cannot parse payload")
	}

	return terms.([]wsdot.Terminal)
}
