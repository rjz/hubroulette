package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type TeamMember struct {
	Github *string `json:"github,omitempty"`
	Slack  *string `json:"slack,omitempty"`
}

// SlackHandle returns the TeamMember's slack handle if available and falls
// back on github if it isn't.
func (tm *TeamMember) SlackHandle() string {
	if tm.Slack != nil {
		return fmt.Sprintf("@%s", *tm.Slack)
	}
	return *tm.Github
}

type TeamMembers []TeamMember

func (users *TeamMembers) FindByGithubLogin(login *string) *TeamMember {
	for _, u := range *users {
		if strings.EqualFold(*u.Github, *login) {
			return &u
		}
	}
	return nil
}

func (users *TeamMembers) Without(login *string) *TeamMembers {
	rest := TeamMembers{}
	for i, u := range *users {
		if !strings.EqualFold(*u.Github, *login) {
			rest = append(rest, (*users)[i])
		}
	}
	return &rest
}

func (users *TeamMembers) Sample() *TeamMember {
	count := len(*users)
	if count == 0 {
		return nil
	}
	t := time.Now()
	rand.Seed(int64(t.Nanosecond()))
	return &(*users)[rand.Intn(count)]
}
