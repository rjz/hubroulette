package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
	client *github.Client
)

func bogusEvent() *github.PullRequestEvent {
	evt := new(github.PullRequestEvent)
	json.Unmarshal([]byte(`{
		"repository": {
			"name": "dingus",
			"owner": { "login":"rjz" }
		},
		"pull_request": {
			"number": 1234,
			"head": { "sha": "abc123"	}
		}
	}`), &evt)
	return evt
}

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = github.NewClient(nil)
	url, _ := url.Parse(server.URL)
	url.Path = "/"
	client.BaseURL = url
	client.UploadURL = url
}

func teardown() {
	server.Close()
}

func TestGetContext(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/rjz/dingus/issues/1234", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"assignee":{"login":"rjz"}}`)
	})

	mux.HandleFunc("/repos/rjz/dingus/git/trees/abc123", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{
			"sha": "abc123",
			"tree":[
				{
					"path": ".hubrouletterc",
					"type": "blob",
					"sha": "def456"
				}
			]
		}`)
	})

	mux.HandleFunc("/repos/rjz/dingus/git/blobs/def456", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"content":"{\"slackChannel\":\"#maybe\",\"team\":[]}","encoding":"utf-8"}`)
	})

	pre, err := GetContext(client, bogusEvent())
	if err != nil {
		t.Fatal(err)
	}

	if *pre.Options.SlackChannel != "#maybe" {
		t.Error(fmt.Sprintf("expected %s, saw %s", "#maybe", *pre.Options.SlackChannel))
	}

	if *pre.Issue.Assignee.Login != "rjz" {
		t.Error(fmt.Sprintf("expected %s, saw %s", "rjz", *pre.Issue.Assignee.Login))
	}
}

func TestGetContextBase64(t *testing.T) {
	// $ export FIXTURE='{"team":[{"github":"rjz"}]}'
	// $ echo "$FIXTURE" | base64
	fixture := "eyJ0ZWFtIjpbeyJnaXRodWIiOiJyanoifV19Cg=="
	// $ echo "$FIXTURE" | wc -c
	fixtureLen := 28

	setup()
	defer teardown()

	mux.HandleFunc("/repos/rjz/dingus/issues/1234", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"assignee":{"login":"rjz"}}`)
	})

	mux.HandleFunc("/repos/rjz/dingus/git/trees/abc123", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{
			"sha": "abc123",
			"tree":[
				{
					"path": ".hubrouletterc",
					"type": "blob",
					"sha": "def456"
				}
			]
		}`)
	})

	mux.HandleFunc("/repos/rjz/dingus/git/blobs/def456", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, fmt.Sprintf(`{"content":"%s","size":%d,"encoding":"base64"}`, fixture, fixtureLen))
	})

	pre, err := GetContext(client, bogusEvent())
	if err != nil {
		t.Fatal(err)
	}

	user := (*pre.Options.TeamMembers)[0]

	if *user.Github != "rjz" {
		t.Error(fmt.Sprintf("expected %s, saw %s", "rjz", *user.Github))
	}
}

func TestGetContextMissingConfig(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/rjz/dingus/issues/1234", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"assignee":{"login":"rjz"}}`)
	})

	mux.HandleFunc("/repos/rjz/dingus/git/trees/abc123", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{
			"sha": "abc123",
			"tree":[]
		}`)
	})

	if _, err := GetContext(client, bogusEvent()); err == nil {
		t.Error(err)
	}
}

func TestGetContextMissingAssignee(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/rjz/dingus/issues/1234", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"assignee":{"login":"rjz"}}`)
	})

	mux.HandleFunc("/repos/rjz/dingus/git/trees/abc123", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{
			"sha": "abc123",
			"tree":[
				{
					"path": ".hubrouletterc",
					"type": "blob",
					"sha": "def456"
				}
			]
		}`)
	})

	if _, err := GetContext(client, bogusEvent()); err == nil {
		t.Error(err)
	}
}

func user(login string) *github.User {
	user := github.User{Login: &login}
	return &user
}

func issueAssignedTo(assignee *github.User) *github.Issue {
	issue := github.Issue{Assignee: assignee}
	return &issue
}

func expectString(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Error(fmt.Sprintf("expected '%s', got %s.", expected, actual))
	}
}

func TestContextIsAssignedFalse(t *testing.T) {
	ctx := PullRequestEventContext{Issue: issueAssignedTo(nil)}
	if ctx.IsAssigned() {
		t.Error("expected user to be assigned, but apparently not.")
	}
}

func TestContextIsAssignedTrue(t *testing.T) {
	ctx := PullRequestEventContext{Issue: issueAssignedTo(user("rjz"))}
	if !ctx.IsAssigned() {
		t.Error("expected user to be assigned, but apparently not.")
	}
}

func TestContextAssigneeSlackHandle(t *testing.T) {
	opts := Options{TeamMembers: someTeamMembers()}
	ctx := PullRequestEventContext{Issue: issueAssignedTo(user("bob")), Options: &opts}
	expectString(t, "bob", ctx.AssigneeSlackHandle())
}

func TestContextAssigneeSlackHandleWithSlack(t *testing.T) {
	opts := Options{TeamMembers: someTeamMembers()}
	ctx := PullRequestEventContext{Issue: issueAssignedTo(user("alice")), Options: &opts}
	expectString(t, "@slack_alice", ctx.AssigneeSlackHandle())
}
