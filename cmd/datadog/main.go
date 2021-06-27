package main

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	file, err := os.OpenFile("ferry.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("Cannot launch logrus")
	}
	log.Out = file

	for {
		time.Sleep(1000 * time.Millisecond)
		log.Info("Heartbeat")
	}
}
