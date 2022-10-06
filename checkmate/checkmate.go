package checkmate

import (
	"context"
	"errors"
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

	action.Debugf("Checklists: %v", checklists)

	allChecked := true
	for _, checklist := range checklists {
		allChecked = allChecked && checklist.AllChecked()

		if !checklist.AllChecked() {
			headerNoPrefix := strings.TrimPrefix(strings.TrimSpace(checklist.Header), "#")
			action.Errorf("Checklist not completed %s", headerNoPrefix)
		}
		action.AddStepSummary(checklist.Summary())
	}

	if !allChecked {
		return errors.New("not all checklists are completed")
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
