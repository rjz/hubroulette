package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/nlopes/slack"
	"golang.org/x/oauth2"
	"log"
	"os"
)

func githubClient(accessToken string) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	oc := oauth2.NewClient(oauth2.NoContext, ts)
	return github.NewClient(oc)
}

func pullRequestEventAttachment(ctx *PullRequestEventContext) *slack.Attachment {
	title := fmt.Sprintf("%s: “%s”", ctx.FullName(), *ctx.PullRequest.Title)
	titleLink := *ctx.PullRequest.HTMLURL
	return &slack.Attachment{
		Title:     title,
		TitleLink: titleLink,
		Text:      *ctx.PullRequest.Body,
		Fallback:  fmt.Sprintf("%s - %s", title, titleLink),
	}
}

func slackMessage(client *slack.Client, ctx *PullRequestEventContext, msg string) error {
	messageParams := slack.PostMessageParameters{
		IconEmoji:   ":game_die:",
		Username:    "Assignee Bot",
		LinkNames:   1,
		Attachments: []slack.Attachment{*pullRequestEventAttachment(ctx)},
	}
	_, _, err := client.PostMessage(*ctx.Options.SlackChannel, msg, messageParams)
	return err
}

// Assign the pull request to a random team member
func Assign(client *github.Client, ctx *PullRequestEventContext) (*TeamMember, error) {
	assignee := ctx.Options.TeamMembers.Without(ctx.PullRequest.User.Login).Sample()
	if assignee == nil {
		return nil, errors.New("No contributors available")
	}

	patch := &github.IssueRequest{Assignee: assignee.Github}
	if _, _, err := client.Issues.Edit(context.TODO(), *ctx.Repository.Owner.Login, *ctx.Repository.Name, *ctx.PullRequest.Number, patch); err != nil {
		return nil, err
	}

	return assignee, nil
}

func HandlePullRequestEvent(evt *github.PullRequestEvent) error {

	if *evt.Action != "opened" {
		// ignore all other events for now
		return nil
	}

	c := NewConfig()
	ghClient := githubClient(c.githubAccessToken)
	slackClient := slack.New(c.slackToken)

	loggerPrefix := fmt.Sprintf("[pr:%s#%d] ", *evt.Repo.FullName, *evt.PullRequest.Number)
	logger := log.New(os.Stdout, loggerPrefix, log.LstdFlags)

	ctx, contextErr := GetContext(ghClient, evt)
	if contextErr != nil {
		logger.Println("Failed retrieving issue context", contextErr)
		return contextErr
	}

	if ctx.IsAssigned() {
		skipMsg := fmt.Sprintf("%s assigned %s to %s", ctx.Author(), ctx.AssigneeSlackHandle(), ctx.FullName())
		if err := slackMessage(slackClient, ctx, skipMsg); err != nil {
			logger.Println("Failed notifying slack", err)
			return err
		}
		return nil
	}

	assignee, assignErr := Assign(ghClient, ctx)
	if assignErr != nil {
		logger.Println("Failed assigning to github", assignErr)
		return assignErr
	}

	assignMsg := fmt.Sprintf("%s opens %s, and %s draws the lucky straw!", ctx.Author(), ctx.FullName(), assignee.SlackHandle())
	if err := slackMessage(slackClient, ctx, assignMsg); err != nil {
		logger.Println("Failed notifying slack", err)
		return err
	}
	return nil
}
