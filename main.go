package main

import (
	"context"
	"fmt"

	"github.com/sethvargo/go-githubactions"

	"github.com/roryq/checkmate/checkmate"
)

func run() error {
	ctx := context.Background()
	action := githubactions.New()

	cfg, err := checkmate.ConfigFromInputs(action)
	if err != nil {
		return err
	}

	fmt.Println(cfg.Checklists)

	return checkmate.Run(ctx, cfg, action)
}

func main() {
	err := run()
	if err != nil {
		githubactions.Fatalf("%v", err)
	}
}
