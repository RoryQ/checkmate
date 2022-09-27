package checkmate

import (
	"context"
	"strings"

	"github.com/sethvargo/go-githubactions"
)

func Run(ctx context.Context, cfg *Config, action *githubactions.Action) error {
	descriptionPR := ""

	for _, checklist := range Parse(descriptionPR) {
		if !checklist.AllChecked() {
			headerNoPrefix := strings.TrimPrefix(strings.TrimSpace(checklist.Header), "#")
			action.Errorf("Checklist not completed %s", headerNoPrefix)
		}
		action.AddStepSummary(checklist.Summary())
	}

	return nil
}
