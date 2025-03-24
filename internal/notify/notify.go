package notify

import (
	"log"
	"time"
)

var lastNotification time.Time
var config = struct {
	Enable   bool
	Interval int
}{
	Enable:   true,
	Interval: 30,
}

func SendNotification(message string) {
	if !config.Enable {
		return
	}

	now := time.Now()
	if now.Sub(lastNotification).Seconds() < float64(config.Interval) {
		return
	}

	log.Printf("提醒: %s", message)
	lastNotification = now
}
