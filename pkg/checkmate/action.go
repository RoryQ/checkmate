package checkmate

import (
	"context"
	"errors"
	"strings"

	"github.com/google/go-github/v48/github"
	"github.com/sethvargo/go-githubactions"

	pr2 "github.com/roryq/checkmate/pkg/pullrequest"
)

func Run(ctx context.Context, cfg *Config, action *githubactions.Action, gh *github.Client) error {
	githubContext, err := action.Context()
	if err != nil {
		return err
	}

	pr, err := pr2.NewClient(action, gh)
	if err != nil {
		return err
	}

	checklists := []Checklist{}
	if len(cfg.PathsChecklists) > 0 {
		action.Infof("Checking changeset for configured paths")
		comment, err := commenter(ctx, *cfg, action, pr)
		if err != nil {
			return err
		}

		if comment != "" {
			action.Infof("Comment checklist %s", comment)
			checklists = Parse(comment)
		}
	}

	descriptionPR, err := getPullRequestBody(githubContext)
	if err != nil {
		return err
	}

	action.Infof("PR Body: %s", descriptionPR)

	checklists = append(Parse(descriptionPR), checklists...)
	return inspect(checklists, action)
}

func inspect(checklists []Checklist, action *githubactions.Action) error {
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
