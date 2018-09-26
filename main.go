package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/rjz/githubhook"
	"log"
	"net/http"
)

const OK_RESPONSE = ":^)"

type WebhookServer struct {
	secret string
}

func String(s string) *string {
	str := new(string)
	*str = s
	return str
}

func Int(i int) *int {
	j := new(int)
	*j = i
	return j
}

func parsePullRequestEvent(hook *githubhook.Hook) *github.PullRequestEvent {
	evt := github.PullRequestEvent{}
	if err := json.Unmarshal(hook.Payload, &evt); false || err != nil {
		// Invalid JSON. Is the webhook configured to send a `application/json` payload?
		// See: https://developer.github.com/webhooks/creating/#content-type
		log.Fatalf("Invalid JSON")
	}
	return &evt
}

func (s *WebhookServer) handle(w http.ResponseWriter, r *http.Request) {

	body, err := githubhook.Parse([]byte(s.secret), r)

	if err != nil {
		log.Println("Webhook parsing failed", err)
		http.Error(w, err.Error(), 400)
		return
	}

	go func() {
		switch body.Event {
		case "pull_request":
			go HandlePullRequestEvent(parsePullRequestEvent(body))
		default:
			log.Printf("Saw unknown event '%s' and waved as it went by.", body.Event)
		}
	}()

	w.Write([]byte(OK_RESPONSE))
}

func main() {

	c := NewConfig()

	server := WebhookServer{
		secret: c.githubWebhookSecret,
	}

	serverPort := fmt.Sprintf(":%d", c.port)

	log.Printf("Listening on %s", serverPort)

	http.HandleFunc("/hooks/github", server.handle)
	http.ListenAndServe(serverPort, nil)
}
