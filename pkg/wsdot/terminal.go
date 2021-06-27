package wsdot

import (
	"encoding/json"
	"fmt"
	"time"
)

type TerminalsEndpoint struct{}

func (e TerminalsEndpoint) url(date time.Time) string {
	day := date.Format("2006-01-02")
	return fmt.Sprintf("/terminals/%s?apiaccesscode=%s", day, apiKey)
}

func (e TerminalsEndpoint) Payload(args EndpointArguments) (string, error) {
	res, err := request(e.url(args.Date))
	return string(res), err
}

func (e TerminalsEndpoint) Unmarshal(data string) (FerryData, error) {
	terminals := make([]Terminal, 0)
	json.Unmarshal([]byte(data), &terminals)
	return terminals, nil

	/*
		terminalsByDesc := make(map[string]Terminal)
		for _, terminal := range terminals {
			terminalsByDesc[terminal.Description] = terminal
		}
		return terminalsByDesc, nil
	*/
}

type Terminal struct {
	ID          int    `json:"TerminalID"`
	Description string `json:"Description"`
}
