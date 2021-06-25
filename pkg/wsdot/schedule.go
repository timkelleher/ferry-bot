package wsdot

import (
	"encoding/json"
	"fmt"
	"time"
)

type ScheduleEndpoint struct{}

func (e ScheduleEndpoint) url(date time.Time, from, to int) string {
	day := date.Format("2006-01-02")
	return fmt.Sprintf("schedule/%s/%d/%d?apiaccesscode=%s",
		day,
		from,
		to,
		apiKey,
	)
}

func (e ScheduleEndpoint) Payload(args EndpointArguments) (string, error) {
	res, err := request(e.url(args.Date, args.Schedule.FromID, args.Schedule.ToID))
	return string(res), err
}

func (e ScheduleEndpoint) Unmarshal(data string) (FerryData, error) {
	schedule := Schedule{}
	json.Unmarshal([]byte(data), &schedule)
	return schedule, nil
}

type ScheduleTime struct {
	VesselID          int    `json:"VesselID"`
	VesselName        string `json:"VesselName"`
	VesselPositionNum int    `json:"VesselPositionNum"`
	DepartingTime     string `json:"DepartingTime"`
	Routes            []int  `json:"Routes"`
}

func (st ScheduleTime) GetDepartingTime() time.Time {
	return WSDOTDate(st.DepartingTime)
}

type TerminalCombo struct {
	ArrivingTerminalID    int            `json:"ArrivingTerminalID"`
	ArrivingTerminalName  string         `json:"ArrivingTerminalName"`
	DepartingTerminalID   int            `json:"DepartingTerminalID"`
	DepartingTerminalName string         `json:"DepartingTerminalName"`
	Times                 []ScheduleTime `json:"Times"`
}

func (tc TerminalCombo) DepartingTimes() []string {
	var times []string
	for _, schTime := range tc.Times {
		times = append(times, schTime.GetDepartingTime().Format("3:04 PM"))
	}
	return times
}

type Schedule struct {
	ID             int             `json:"ScheduleID"`
	Name           string          `json:"ScheduleName"`
	Start          string          `json:"ScheduleStart"`
	End            string          `json:"ScheduleEnd"`
	TerminalCombos []TerminalCombo `json:"TerminalCombos"`
}
