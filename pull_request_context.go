package main

import (
	"errors"
	"fmt"
	"github.com/google/go-github/github"
)

const configFileName = ".hubrouletterc"

type PullRequestEventContext struct {
	Repository  *github.Repository
	PullRequest *github.PullRequest
	Issue       *github.Issue
	Options     *Options
}

func (ctx *PullRequestEventContext) Author() string {
	return *ctx.PullRequest.User.Login
}

func (ctx *PullRequestEventContext) FullName() string {
	return fmt.Sprintf("%s#%d", *ctx.Repository.FullName, *ctx.PullRequest.Number)
}

// Look for the SHA of the config file in project root
func configShaFromTree(tree *github.Tree) *string {
	var blobSha *string
	for _, entry := range tree.Entries {
		if *entry.Path == configFileName {
			blobSha = entry.SHA
		}
	}
	return blobSha
}

// Retrieve configuration from the repo, if checked in
func getRc(client *github.Client, owner, repo, sha string) (*Options, error) {

	tree, _, treeErr := client.Git.GetTree(owner, repo, sha, false)
	if treeErr != nil {
		return nil, treeErr
	}

	blobSha := configShaFromTree(tree)
	if blobSha == nil {
		return nil, nil
	}

	blob, _, blobErr := client.Git.GetBlob(owner, repo, *blobSha)
	if blobErr != nil {
		return nil, blobErr
	}

	bytes := []byte(*blob.Content)

	opts, optsErr := ParseOptions(&bytes)
	if optsErr != nil {
		return nil, optsErr
	}

	return opts, nil
}

// Retrieve the context around a pull request event
func GetContext(client *github.Client, evt *github.PullRequestEvent) (*PullRequestEventContext, error) {
	opts, optsErr := getRc(client, *evt.Repo.Owner.Login, *evt.Repo.Name, *evt.PullRequest.Head.SHA)
	if optsErr != nil {
		return nil, optsErr
	}

	if opts == nil {
		// compat.ab: allow globally config'd rc for now
		c := NewConfig()
		if c.hubrouletterc == nil {
			return nil, errors.New("No options available")
		}
		opts = c.hubrouletterc
		// compat.fin
	}

	issue, _, issueErr := client.Issues.Get(*evt.Repo.Owner.Login, *evt.Repo.Name, *evt.PullRequest.Number)
	if issueErr != nil {
		return nil, issueErr
	}

	pre := PullRequestEventContext{
		Repository:  evt.Repo,
		PullRequest: evt.PullRequest,
		Options:     opts,
		Issue:       issue,
	}

	return &pre, nil
}
