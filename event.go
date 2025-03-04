// event.go
package main

import (
	"fmt"

	"google.golang.org/api/gmail/v1"
)

func triggerCRMEvent(srv *gmail.Service, sender, clientEmail, clientName string) error {
	employeeEmail := getEmployeeEmail()
	data := generatePersonalizedContent(clientName, employeeEmail)
	return sendTemplatedEmail(srv, sender, clientEmail, fmt.Sprintf("Welcome %s!", clientName), data)
}
