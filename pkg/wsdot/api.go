package wsdot

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

var baseURL = "https://www.wsdot.wa.gov/ferries/api/schedule/rest"
var apiKey = "5c72f79c-27bb-4aad-8e50-43603772b9be"

func request(url string) []byte {
	logrus.Debug("API Call Start...")
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Print(err.Error())
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	logrus.Debug(resp.StatusCode)
	return bodyBytes
}
