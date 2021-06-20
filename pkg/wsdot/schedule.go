package wsdot

import (
	"encoding/json"
	"fmt"
	"time"
)

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

func GetSchedule(date time.Time, from, to Terminal) Schedule {
	payload := request(buildScheduleURL(date, from, to))

	schedule := Schedule{}
	json.Unmarshal(payload, &schedule)
	return schedule
}
