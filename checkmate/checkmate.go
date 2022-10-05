package checkmate

import (
	"context"
	"strings"

	"github.com/sethvargo/go-githubactions"
)

func Run(ctx context.Context, cfg *Config, action *githubactions.Action) error {
	descriptionPR, err := getPullRequestBody(action)
	if err != nil {
		return err
	}

	action.Infof("PR Body: %s", descriptionPR)

	checklists := Parse(descriptionPR)

	action.Infof("Checklists: %v", checklists)

	for _, checklist := range checklists {
		if !checklist.AllChecked() {
			headerNoPrefix := strings.TrimPrefix(strings.TrimSpace(checklist.Header), "#")
			action.Errorf("Checklist not completed %s", headerNoPrefix)
		}
		action.AddStepSummary(checklist.Summary())
	}

	return nil
}

func getPullRequestBody(action *githubactions.Action) (string, error) {
	ghctx, err := action.Context()
	if err != nil {
		return "", err
	}
	body := ghctx.Event["pull_request"].(map[string]any)["body"]
	return body.(string), nil
}
