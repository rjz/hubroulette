package main

import (
	"errors"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/nlopes/slack"
	"golang.org/x/oauth2"
	"log"
	"os"
)

type PRHandler struct {
	assignee           *github.User
	githubClient       *github.Client
	githubPullRequest  *github.PullRequest
	githubRepository   *github.Repository
	slackChannel       string
	slackClient        *slack.Client
	slackMessageParams slack.PostMessageParameters
}

func githubClient(accessToken string) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	oc := oauth2.NewClient(oauth2.NoContext, ts)
	return github.NewClient(oc)
}

func HandlePullRequestEvent(evt *github.PullRequestEvent) error {

	c := NewConfig()

	contributors := NewAssignees(c.assigneeLogins)

	if *evt.Action != "opened" {
		// ignore all other events for now
		return nil
	}

	ghClient := githubClient(c.githubAccessToken)

	h := PRHandler{
		assignee:          nil,
		githubClient:      ghClient,
		githubPullRequest: evt.PullRequest,
		githubRepository:  evt.Repo,
		slackClient:       slack.New(c.slackToken),
		slackChannel:      c.slackChannel,
		slackMessageParams: slack.PostMessageParameters{
			IconEmoji: ":game_die:",
		},
	}

	loggerPrefix := fmt.Sprintf("[pr:%s#%d] ", *evt.Repo.FullName, *evt.PullRequest.Number)
	logger := log.New(os.Stdout, loggerPrefix, log.LstdFlags)

	if err := h.currentAssignee(*evt.PullRequest.Number); err != nil {
		logger.Println("Failed retrieving issue status", err)
		return err
	}

	if h.assignee != nil {
		logger.Println("Issue is already assigned, skipping")
		return nil
	}

	if err := h.Assign(contributors.Without(evt.PullRequest.User.Login).Sample()); err != nil {
		logger.Println("Failed assigning to github", err)
		return err
	}

	if err := h.Notify(); err != nil {
		logger.Println("Failed notifying slack", c.slackChannel, err)
		return err
	}

	return nil
}

// Look up current assignee for github issue and set `h.assignee`
func (h *PRHandler) currentAssignee(number int) error {

	issue, _, err := h.githubClient.Issues.Get(*h.githubRepository.Owner.Login, *h.githubRepository.Name, number)

	if err != nil {
		return err
	}

	h.assignee = issue.Assignee

	return nil
}

// Notify slack about a new assignee
func (h *PRHandler) Notify() error {
	prName := fmt.Sprintf("%s#%d", *h.githubRepository.FullName, *h.githubPullRequest.Number)
	msg := fmt.Sprintf("%s is now open. %s drew the lucky straw!\n %s", prName, *h.assignee.Login, *h.githubPullRequest.HTMLURL)

	_, _, err := h.slackClient.PostMessage(h.slackChannel, msg, h.slackMessageParams)
	return err
}

// Assign the pull request to the specified `assignee`
func (h *PRHandler) Assign(assignee *github.User) error {

	if assignee == nil {
		return errors.New("No contributors available")
	}

	patch := &github.IssueRequest{Assignee: assignee.Login}

	if _, _, err := h.githubClient.Issues.Edit(*h.githubRepository.Owner.Login, *h.githubRepository.Name, *h.githubPullRequest.Number, patch); err != nil {
		return err
	}

	h.assignee = assignee

	return nil
}
