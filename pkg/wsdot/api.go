package wsdot

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

var baseURL = "https://www.wsdot.wa.gov/ferries/api/schedule/rest"
var apiKey = os.Getenv("WSDOT_API_KEY")

var validWSDOTDateTime = regexp.MustCompile(`(\d+)`)

var ErrMissingArg = errors.New("missing required argument")
var ErrInvalidArg = errors.New("invalid argument")
var ErrResponseStatusCode = errors.New("invalid status code")

type Endpoint interface {
	Payload(EndpointArguments) (string, error)
	Unmarshal(string) (FerryData, error)
}

type EndpointArguments struct {
	Date     time.Time
	Schedule struct {
		FromID int
		ToID   int
	}
}

type FerryData interface{}

const (
	CacheFlushDateID = "cache_flush_date"
	TerminalsID      = "terminals"
	ScheduleID       = "schedule"
)

var Endpoints = map[string]Endpoint{
	CacheFlushDateID: CacheFlushDateEndpoint{},
	TerminalsID:      TerminalsEndpoint{},
	ScheduleID:       ScheduleEndpoint{},
}

func request(endpoint string) ([]byte, error) {
	url := baseURL + "/" + endpoint
	fmt.Println("API Call Start...", url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	fmt.Println("API Status Code:", resp.StatusCode)
	if resp.StatusCode != 200 {
		return bodyBytes, ErrResponseStatusCode
	}

	return bodyBytes, nil
}

func WSDOTDate(date string) time.Time {
	timestamp, err := strconv.Atoi(validWSDOTDateTime.FindString(date))
	if err != nil {
		return time.Now()
	}

	return time.Unix(int64(timestamp/1000), 0)
}
