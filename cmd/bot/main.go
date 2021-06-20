package main

import (
	"time"

	"github.com/sanity-io/litter"
	"github.com/timkelleher/ferry-bot/pkg/wsdot"
)

func main() {
	terminals := wsdot.Terminals(time.Now())
	litter.Dump(terminals)

	schedule := wsdot.ScheduleTF(time.Now(), terminals["Bainbridge Island"], terminals["Seattle"])
	litter.Dump(schedule)
}
