package main

import (
	"encoding/json"
	"errors"
)

type Options struct {
	SlackChannel *string      `json:"slackChannel,omitempty"`
	TeamMembers  *TeamMembers `json:"team,omitempty"`
}

func ParseOptions(rawOpts *[]byte) (*Options, error) {
	o := &Options{
		SlackChannel: String("#general"),
	}

	if err := json.Unmarshal(*rawOpts, o); err != nil {
		return nil, err
	}

	if o.TeamMembers == nil {
		return nil, errors.New(".team is required")
	}

	return o, nil
}
