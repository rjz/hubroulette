package main

import (
	"testing"
)

func someAssignees() []string {
	return []string{
		"alice",
		"bob",
		"mallory",
	}
}

func TestNewAssignees(t *testing.T) {
	assignees := NewAssignees(someAssignees())
	count := len(*assignees)
	if count != 3 {
		t.Errorf("Expected 3, got %d", count)
	}
}

func TestAssigneesWithout(t *testing.T) {
	assignees := NewAssignees(someAssignees())
	count := len(*assignees)
	expectedCount := count - 1
	others := assignees.Without(&someAssignees()[0])

	if len(*others) != expectedCount {
		t.Errorf("Expected %d, got %d", expectedCount, len(*others))
	}
}

func TestSampleUsers(t *testing.T) {
	assignee := NewAssignees(someAssignees()).Sample()
	if assignee == nil {
		t.Errorf("Expected user, got nil")
	}

	login := *assignee.Login

	if login != "bob" && login != "alice" && login != "mallory" {
		t.Errorf("Expected someone familiar, didn't get 'em")
	}
}
