package main

import (
	"log"
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
	port                int64
	hubrouletterc       *Options
}

func NewConfig() *Config {

	var err error

	logger := log.New(os.Stdout, "[config] ", log.LstdFlags)

	config := Config{
		assigneeLogins:      strings.Split(os.Getenv("ASSIGNEE_LOGINS"), ","),
		githubAccessToken:   os.Getenv("GITHUB_ACCESS_TOKEN"),
		githubWebhookSecret: os.Getenv("GITHUB_WEBHOOK_SECRET"),
		slackUsername:       os.Getenv("SLACK_USERNAME"),
		slackChannel:        os.Getenv("SLACK_CHANNEL"),
		slackToken:          os.Getenv("SLACK_TOKEN"),
	}

	port := os.Getenv("PORT")
	if config.port, err = strconv.ParseInt(port, 10, 32); err != nil {
		logger.Printf("Environment specified invalid PORT '%s', using default (%d)", port, default_port)
		config.port = default_port
	}

	hrc := []byte(os.Getenv("HUBROULETTERC"))
	if config.hubrouletterc, err = ParseOptions(&hrc); err != nil {
		logger.Printf("Environment specified invalid HUBROULETTERC '%s'", string(hrc))
	}

	return &config
}
