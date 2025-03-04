// logger.go
package main

import (
	"fmt"
	"time"
)

var emailLogs []EmailLog

type EmailLog struct {
	Timestamp   time.Time
	ToEmail     string
	Subject     string
	SenderEmail string
	Status      string
}

func displayEmailTimeline() {
	for _, log := range emailLogs {
		fmt.Printf("[%s] To: %-20s Subject: %-20s Status: %s\n",
			log.Timestamp.Format("2006-01-02 15:04"),
			log.ToEmail,
			log.Subject,
			log.Status,
		)
	}
}
