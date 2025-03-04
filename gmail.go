// gmail.go
package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func initGmailService() (*gmail.Service, string) {
	config := &oauth2.Config{
		ClientID:     os.Getenv("GMAIL_CLIENT_ID"),
		ClientSecret: os.Getenv("GMAIL_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GMAIL_REDIRECT_URI"),
		Scopes:       []string{gmail.GmailSendScope},
		Endpoint:     google.Endpoint,
	}

	client := getClient(config)
	srv, err := gmail.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatal(err)
	}

	user, err := srv.Users.GetProfile("me").Do()
	if err != nil {
		log.Fatal(err)
	}

	return srv, user.EmailAddress
}

func getClient(config *oauth2.Config) *http.Client {
	tok, err := tokenFromFile("token.json")
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken("token.json", tok)
	}
	return config.Client(context.Background(), tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to auth URL and enter code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}

	tok, err := config.Exchange(context.TODO(), code)
	if err != nil {
		log.Fatal(err)
	}
	return tok
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func sendTemplatedEmail(srv *gmail.Service, sender, to, subject string, data map[string]string) error {
	const tpl = `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}
MIME-Version: 1.0
Content-Type: text/plain; charset="utf-8"

Dear {{.ClientName}},
Welcome! Your contact is {{.EmployeeName}}.
`

	tmpl, err := template.New("email").Parse(tpl)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, struct {
		From, To, Subject, ClientName, EmployeeName string
	}{
		sender,
		to,
		subject,
		data["ClientName"],
		data["EmployeeName"],
	}); err != nil {
		return err
	}

	msg := gmail.Message{
		Raw: base64.URLEncoding.EncodeToString(body.Bytes()),
	}

	_, err = srv.Users.Messages.Send("me", &msg).Do()
	if err != nil {
		return err
	}

	emailLogs = append(emailLogs, EmailLog{
		Timestamp:   time.Now(),
		ToEmail:     to,
		Subject:     subject,
		SenderEmail: sender,
		Status:      "sent",
	})

	return nil
}
