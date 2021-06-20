package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/timkelleher/ferry-bot/pkg/wsdot"
)

var discordWebhookUrl = os.Getenv("DISCORD_WEBHOOK_URL")

func formatTimes(times []string) string {
	return strings.Join(times, " | ")
}

func discordWebhookPing() error {
	fmt.Println("Pinging Discord Webhook server...")
	_, err := http.Get(discordWebhookUrl)
	if err != nil {
		return err
	}
	return nil
}

func discordWebhookRequest(title, times string) error {
	requestBody, err := json.Marshal(map[string]string{
		"title":                 title,
		"times":                 times,
		"discord_webhook_id":    os.Getenv("DISCORD_WEBHOOK_ID"),
		"discord_webhook_token": os.Getenv("DISCORD_WEBHOOK_TOKEN"),
	})
	if err != nil {
		return err
	}

	url := discordWebhookUrl + "/ferry"

	fmt.Println("Sending request to Discord Webhook server...")
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := discordWebhookPing()
	if err != nil {
		panic(err)
	}

	date := time.Now()
	terminals := wsdot.GetTerminals(date)

	title := fmt.Sprintf("%s - Bainbridge Island -> Seattle\n", date.Format("Mon Jan _2"))
	scheduleBI2S := wsdot.GetSchedule(date, terminals["Bainbridge Island"], terminals["Seattle"])
	err = discordWebhookRequest(title, formatTimes(scheduleBI2S.TerminalCombos[0].DepartingTimes()))
	if err != nil {
		panic(err)
	}

	title = fmt.Sprintf("%s - Seattle -> Bainbridge Island\n", date.Format("Mon Jan _2"))
	scheduleS2BI := wsdot.GetSchedule(date, terminals["Seattle"], terminals["Bainbridge Island"])
	err = discordWebhookRequest(title, formatTimes(scheduleS2BI.TerminalCombos[0].DepartingTimes()))
	if err != nil {
		panic(err)
	}
}
