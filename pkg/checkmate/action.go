package checkmate

import (
	"context"
	"errors"
	"strings"

	"github.com/google/go-github/v48/github"
	"github.com/sethvargo/go-githubactions"
)

func Run(ctx context.Context, cfg *Config, action *githubactions.Action, gh *github.Client) error {
	githubContext, err := action.Context()
	if err != nil {
		return err
	}

	if len(cfg.PathsChecklists) > 0 {
		action.Infof("Checking changeset for configured paths")
		comment, err := commenter(ctx, *cfg, action, gh)
		if err != nil {
			return err
		}

		if comment != "" {
			action.Infof("Comment checklist %s", comment)
			if err := inspect(comment, action); err != nil {
				return err
			}
		}
	}

	descriptionPR, err := getPullRequestBody(githubContext)
	if err != nil {
		return err
	}

	action.Infof("PR Body: %s", descriptionPR)

	return inspect(descriptionPR, action)
}

func inspect(body string, action *githubactions.Action) error {
	checklists := Parse(body)

	action.Debugf("Checklists: %v", checklists)

	action.AddStepSummary("_The following checklists were found and validated:_\n")

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

func getPullRequestBody(ghctx *githubactions.GitHubContext) (string, error) {
	body, ok := ghctx.Event["pull_request"].(map[string]any)["body"]
	if !ok || body == nil {
		return "", nil
	}
	return body.(string), nil
}
