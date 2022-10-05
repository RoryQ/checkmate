package main

import (
	"context"

	"github.com/sethvargo/go-githubactions"

	"github.com/roryq/checkmate/checkmate"
)

func run() error {
	ctx := context.Background()
	action := githubactions.New()

	return checkmate.Run(ctx, &checkmate.Config{}, action)
}

func main() {
	err := run()
	if err != nil {
		githubactions.Fatalf("%v", err)
	}
}
