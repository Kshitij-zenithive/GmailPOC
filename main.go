// main.go
package main

import (
	"log"
)

func main() {
	// Initialize Gmail service
	srv, senderEmail := initGmailService()

	// Simulate CRM events
	clients := []struct {
		email string
		name  string
	}{
		{"client1@example.com", "Alice"},
		{"client2@example.com", "Bob"},
	}

	for _, client := range clients {
		if err := triggerCRMEvent(srv, senderEmail, client.email, client.name); err != nil {
			log.Printf("Failed to process %s: %v", client.email, err)
		}
	}

	// Display email timeline
	displayEmailTimeline()
}
