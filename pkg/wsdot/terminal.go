package wsdot

import (
	"encoding/json"
	"fmt"
	"time"
)

type Terminal struct {
	ID          int    `json:"TerminalID"`
	Description string `json:"Description"`
}

func buildTerminalsURL(date time.Time) string {
	//day := time.Now().AddDate(0, 0, 5).Format("2006-01-02")
	day := date.Format("2006-01-02")
	return fmt.Sprintf("%s/terminals/%s?apiaccesscode=%s", baseURL, day, apiKey)
}

func GetTerminals(date time.Time) map[string]Terminal {
	payload := request(buildTerminalsURL(date))

	terminals := make([]Terminal, 0)
	json.Unmarshal(payload, &terminals)

	terminalsByDesc := make(map[string]Terminal)
	for _, terminal := range terminals {
		terminalsByDesc[terminal.Description] = terminal
	}

	return terminalsByDesc
}
