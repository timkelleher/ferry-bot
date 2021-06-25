package wsdot

import (
	"time"
)

type CacheFlushDateEndpoint struct{}

func (e CacheFlushDateEndpoint) url() string {
	return "cacheflushdate"
}

func (e CacheFlushDateEndpoint) Payload(args EndpointArguments) (string, error) {
	res, err := request(e.url())
	return string(res), err
}

func (e CacheFlushDateEndpoint) Unmarshal(data string) (FerryData, error) {
	return CacheFlushDate{date: data}, nil
}

type CacheFlushDate struct {
	date string `json:"Date"`
}

func (cfd CacheFlushDate) Date() time.Time {
	return WSDOTDate(cfd.date)
}

func (cfd CacheFlushDate) Until() time.Duration {
	return time.Until(cfd.Date())
}
