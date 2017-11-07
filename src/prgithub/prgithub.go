package prgithub

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

//client is the github API client
var client = &github.Client{}

//ctx is the github API context
var ctx context.Context

//New initializes the client, a API token has to be provided
func New(token string) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	ctx = context.Background()
	oa := oauth2.NewClient(ctx, ts)
	client = github.NewClient(oa)
}

//Get returns the github client and the github ctx
func Get() (context.Context, *github.Client) {
	return ctx, client
}
