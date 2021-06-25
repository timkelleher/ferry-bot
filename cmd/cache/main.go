package main

import (
	"context"
	"fmt"
	"time"

	"github.com/sanity-io/litter"
	"github.com/timkelleher/ferry-bot/pkg/cache"
	"github.com/timkelleher/ferry-bot/pkg/wsdot"
)

func main() {
	ctx := context.Background()

	endpoint := wsdot.Endpoints[wsdot.CacheFlushDateID]
	payload, err := cache.Get(ctx, wsdot.CacheFlushDateID, time.Hour, endpoint.Payload, wsdot.EndpointArguments{})
	if err != nil {
		panic("cannot fetch payload")
	}
	cacheTime, err := endpoint.Unmarshal(payload)
	if err != nil {
		panic("cannot unmarshal payload")
	}

	cacheTimeObj := cacheTime.(wsdot.CacheFlushDate)
	fmt.Println(cacheTimeObj.Until())

	endpoint = wsdot.Endpoints[wsdot.TerminalsID]
	payload, err = cache.Get(ctx, wsdot.TerminalsID, 24*time.Hour, endpoint.Payload, wsdot.EndpointArguments{
		Date: time.Now(),
	})
	if err != nil {
		panic("cannot fetch payload")
	}
	terms, err := endpoint.Unmarshal(payload)
	if err != nil {
		panic("cannot unmarshal payload")
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
	payload, err = cache.Get(ctx, wsdot.ScheduleID, time.Hour, endpoint.Payload, args)
	if err != nil {
		panic("cannot fetch payload")
	}
	bainbridgeSchedule, err := endpoint.Unmarshal(payload)
	if err != nil {
		panic("cannot unmarshal payload")
	}

	bainbridgeScheduleObj := bainbridgeSchedule.(wsdot.Schedule)
	litter.Dump(bainbridgeScheduleObj.TerminalCombos[0].DepartingTimes())
}
