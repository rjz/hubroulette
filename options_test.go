package main

import (
	"fmt"
	"testing"
)

func TestParseInvalidOptions(t *testing.T) {
	rawOpts := []byte(`foobar`)
	_, err := ParseOptions(&rawOpts)
	if err == nil {
		t.Error("expected failure; didn't get it.")
	}
}

func TestOptionsMissingTeam(t *testing.T) {
	rawOpts := []byte(`{"slackChannel":"#general"}`)
	_, err := ParseOptions(&rawOpts)
	if err == nil {
		t.Error("expected error for missing .team; didn't get it.")
	}
}

func TestParseValidOptions(t *testing.T) {
	rawOpts := []byte(`{"slackChannel":"#general","team":[]}`)
	opts, _ := ParseOptions(&rawOpts)
	if *opts.SlackChannel != "#general" {
		t.Error(fmt.Sprintf("expected %s, got %s", *opts.SlackChannel, "#general"))
	}
}

func TestParseEmptyTeam(t *testing.T) {
	rawOpts := []byte(`{"slackChannel":"#general","team":[]}`)
	opts, _ := ParseOptions(&rawOpts)
	if opts.TeamMembers == nil {
		t.Error(fmt.Sprintf("expected Team Members"))
	}
}
