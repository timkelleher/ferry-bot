package wsdot

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

var baseURL = "https://www.wsdot.wa.gov/ferries/api/schedule/rest"
var apiKey = os.Getenv("WSDOT_API_KEY")

var validWSDOTDateTime = regexp.MustCompile(`(\d+)`)

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

func WSDOTDate(date string) time.Time {
	timestamp, err := strconv.Atoi(validWSDOTDateTime.FindString(date))
	if err != nil {
		return time.Now()
	}

	return time.Unix(int64(timestamp/1000), 0)
}
