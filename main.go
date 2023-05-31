package main

import (
	"context"
	"net/http"

	"github.com/google/go-github/v48/github"
	"github.com/sethvargo/go-githubactions"
	"golang.org/x/oauth2"

	"github.com/marcusvnac/checkmate-evo/pkg/checkmate"
	"github.com/marcusvnac/checkmate-evo/pkg/checkmate/inputs"
)

func run() error {
	ctx := context.Background()
	action := githubactions.New()

	cfg, err := checkmate.ConfigFromInputs(action)
	if err != nil {
		return err
	}

	var tc *http.Client
	if token := action.GetInput(inputs.GithubToken); token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc = oauth2.NewClient(ctx, ts)
	}

	gh := github.NewClient(tc)
	return checkmate.Run(ctx, cfg, action, gh)
}

func main() {
	err := run()
	if err != nil {
		githubactions.Fatalf("%v", err)
	}
}
