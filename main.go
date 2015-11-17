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

func parsePullRequestEvent(hook *githubhook.Hook) *github.PullRequestEvent {
	evt := github.PullRequestEvent{}
	json.Unmarshal(hook.Payload, &evt)
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
