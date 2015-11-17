package main

import (
	"github.com/google/go-github/github"
	"math/rand"
	"strings"
	"time"
)

type Assignees []github.User

func NewAssignees(assigneeStrings []string) *Assignees {
	assignees := make(Assignees, len(assigneeStrings))

	for i, _ := range assigneeStrings {
		assignees[i] = github.User{Login: &assigneeStrings[i]}
	}

	return &assignees
}

func (as *Assignees) Without(login *string) *Assignees {
	rest := Assignees{}
	for i, u := range *as {
		if !strings.EqualFold(*u.Login, *login) {
			rest = append(rest, (*as)[i])
		}
	}
	return &rest
}

func (as *Assignees) Sample() *github.User {
	count := len(*as)
	if count == 0 {
		return nil
	}
	t := time.Now()
	rand.Seed(int64(t.Nanosecond()))
	return &(*as)[rand.Intn(count)]
}
