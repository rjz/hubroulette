package main

import (
	"fmt"
	"github.com/google/go-github/github"
	"net/http"
	"testing"
)

func bogusContext(rc []byte) *PullRequestEventContext {
	opts, _ := ParseOptions(&rc)
	user := github.User{Login: String("rjz")}
	pr := github.PullRequest{User: &user, Number: Int(1234)}
	repo := github.Repository{Name: String("dingus"), Owner: &user}

	return &PullRequestEventContext{
		Options:     opts,
		PullRequest: &pr,
		Repository:  &repo,
	}
}

func TestAssignEmptyTeam(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/rjz/dingus/issues/1234", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{}`)
	})

	ctx := bogusContext([]byte(`{"team":[]}`))

	if _, err := Assign(client, ctx); err == nil {
		t.Error("Expected no valid assignee")
	}
}

func TestAssignNoTeamMemberAvailable(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/rjz/dingus/issues/1234", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{}`)
	})

	ctx := bogusContext([]byte(`{"team":[{"github":"rjz"}]}`))

	if _, err := Assign(client, ctx); err == nil {
		t.Error("Expected no valid assignee")
	}
}

func TestAssignTeamMemberAvailable(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/rjz/dingus/issues/1234", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{}`)
	})

	ctx := bogusContext([]byte(`{"team":[{"github":"rjz"},{"github":"example"}]}`))

	if assignee, _ := Assign(client, ctx); *assignee.Github != "example" {
		t.Error("Expected assignee, didn't get it.")
	}
}
