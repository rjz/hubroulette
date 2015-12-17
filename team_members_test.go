package main

import (
	"testing"
)

func someTeamMembers() *TeamMembers {
	return &TeamMembers{
		{Github: String("alice"), Slack: String("slack_alice")},
		{Github: String("bob")},
		{Github: String("mallory")},
	}
}

func TestTeamMembersFindByGithubLogin(t *testing.T) {
	alice := someTeamMembers().FindByGithubLogin(String("alice"))
	if *alice.Slack != "slack_alice" {
		t.Errorf("Didn't expect to see '%s'", *alice.Slack)
	}
}

func TestTeamMembersWithout(t *testing.T) {
	others := someTeamMembers().Without(String("bob"))
	if len(*others) != 2 {
		t.Errorf("Expected %d, got %d", 3, len(*others))
	}
}

func TestSampleUsers(t *testing.T) {
	login := *(someTeamMembers().Sample()).Github
	if login != "bob" && login != "alice" && login != "mallory" {
		t.Errorf("Expected someone familiar, didn't get 'em")
	}
}

func TestSampleSingleUser(t *testing.T) {
	teamMember := someTeamMembers().Without(String("bob")).Without(String("mallory")).Sample()
	if *teamMember.Github != "alice" {
		t.Errorf("Expected someone familiar, didn't get 'em")
	}
}
