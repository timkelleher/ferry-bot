package main

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/timkelleher/ferry-bot/pkg/cache"
	"github.com/timkelleher/ferry-bot/pkg/wsdot"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	tracer.Start(
		tracer.WithEnv("prod"),
		tracer.WithService("ferry"),
		tracer.WithServiceVersion("0.1"),
	)
	defer tracer.Stop()

	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	file, err := os.OpenFile("ferry.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("Cannot launch logrus")
	}
	log.Out = file

	endpoint := wsdot.Endpoints[wsdot.CacheFlushDateID]
	payload, err := cache.Get(wsdot.CacheFlushDateID, time.Hour, log, endpoint.Payload, wsdot.EndpointArguments{})
	if err != nil {
		log.WithFields(logrus.Fields{
			"endpoint": wsdot.CacheFlushDateID,
		}).Fatalf("Cannot fetch payload")
	}
	cacheTime, err := endpoint.Unmarshal(payload)
	if err != nil {
		log.WithFields(logrus.Fields{
			"endpoint": wsdot.CacheFlushDateID,
		}).Fatalf("Cannot parse payload")
	}

	cacheTimeObj := cacheTime.(wsdot.CacheFlushDate)
	fmt.Println(cacheTimeObj.Until())
	log.WithFields(logrus.Fields{
		"value": cacheTimeObj.Date(),
		"until": cacheTimeObj.Until(),
	}).Info("WSDOT cache expiration time")

	endpoint = wsdot.Endpoints[wsdot.TerminalsID]
	payload, err = cache.Get(wsdot.TerminalsID, 24*time.Hour, log, endpoint.Payload, wsdot.EndpointArguments{
		Date: time.Now(),
	})
	if err != nil {
		log.WithFields(logrus.Fields{
			"endpoint": wsdot.TerminalsID,
		}).Fatalf("Cannot fetch payload")
	}
	terms, err := endpoint.Unmarshal(payload)
	if err != nil {
		log.WithFields(logrus.Fields{
			"endpoint": wsdot.TerminalsID,
		}).Fatalf("Cannot parse payload")
	}

	terminalsMap := terms.(map[string]wsdot.Terminal)
	for _, terminal := range terminalsMap {
		fmt.Println(terminal.ID, terminal.Description)
	}

	endpoint = wsdot.Endpoints[wsdot.ScheduleID]
	args := wsdot.EndpointArguments{
		Date: time.Now(),
		Schedule: struct {
			FromID int
			ToID   int
		}{terminalsMap["Bainbridge Island"].ID, terminalsMap["Seattle"].ID},
	}
	payload, err = cache.Get(wsdot.ScheduleID, time.Hour, log, endpoint.Payload, args)
	if err != nil {
		log.WithFields(logrus.Fields{
			"endpoint": wsdot.ScheduleID,
		}).Fatalf("Cannot fetch payload")
	}
	bainbridgeSchedule, err := endpoint.Unmarshal(payload)
	if err != nil {
		log.WithFields(logrus.Fields{
			"endpoint": wsdot.ScheduleID,
		}).Fatalf("Cannot parse payload")
	}

	bainbridgeScheduleObj := bainbridgeSchedule.(wsdot.Schedule)
	log.WithFields(logrus.Fields{
		"times": bainbridgeScheduleObj.TerminalCombos[0].DepartingTimes(),
	}).Info("Found departing times")
}
