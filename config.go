package main

import (
	"os"
	"strconv"
	"strings"
)

const default_port = 3000

type Config struct {
	assigneeLogins      []string
	githubAccessToken   string
	githubWebhookSecret string
	slackUsername       string
	slackChannel        string
	slackToken          string
	port                int
}

func NewConfig() *Config {

	port, err := strconv.ParseInt(os.Getenv("PORT"), 10, 32)
	if err != nil {
		port = default_port
	}

	config := Config{
		assigneeLogins:      strings.Split(os.Getenv("ASSIGNEE_LOGINS"), ","),
		githubAccessToken:   os.Getenv("GITHUB_ACCESS_TOKEN"),
		githubWebhookSecret: os.Getenv("GITHUB_WEBHOOK_SECRET"),
		slackUsername:       os.Getenv("SLACK_USERNAME"),
		slackChannel:        os.Getenv("SLACK_CHANNEL"),
		slackToken:          os.Getenv("SLACK_TOKEN"),
		port:                int(port),
	}

	return &config
}
