package main

import (
	"context"

	"github.com/google/go-github/v48/github"
	"github.com/sethvargo/go-githubactions"

	"github.com/roryq/checkmate/pkg/checkmate"
)

func run() error {
	ctx := context.Background()
	action := githubactions.New()

	cfg, err := checkmate.ConfigFromInputs(action)
	if err != nil {
		return err
	}

	gh := github.NewClient(nil)
	return checkmate.Run(ctx, cfg, action, gh)
}

func main() {
	err := run()
	if err != nil {
		githubactions.Fatalf("%v", err)
	}
}
