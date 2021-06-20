package wsdot

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sanity-io/litter"
)

type Schedule struct {
	ID   int    `json:"ScheduleID"`
	Name string `json:"ScheduleName"`
}

func buildScheduleURL(date time.Time, from, to Terminal) string {
	day := date.Format("2006-01-02")
	return fmt.Sprintf("%s/schedule/%s/%d/%d?apiaccesscode=%s",
		baseURL,
		day,
		from.ID,
		to.ID,
		apiKey,
	)
}

func ScheduleTF(date time.Time, from, to Terminal) Schedule {
	litter.Dump(buildScheduleURL(date, from, to))
	payload := request(buildScheduleURL(date, from, to))

	var schedule Schedule
	json.Unmarshal(payload, &schedule)
	return schedule
}
