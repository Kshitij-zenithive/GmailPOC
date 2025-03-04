// ai.go
package main

func generatePersonalizedContent(client, employee string) map[string]string {
	return map[string]string{
		"ClientName":   client,
		"EmployeeName": employee,
	}
}
