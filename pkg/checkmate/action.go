package checkmate

import (
	"context"
	"errors"
	"strings"

	"github.com/google/go-github/v48/github"
	"github.com/sethvargo/go-githubactions"

	"github.com/roryq/checkmate/pkg/pullrequest"
)

func Run(ctx context.Context, cfg *Config, action *githubactions.Action, gh *github.Client) error {
	githubContext, err := action.Context()
	if err != nil {
		return err
	}
	switch githubContext.EventName {
	case "merge_group":
		action.Infof("skipping checkmate on %s", githubContext.EventName)
		return nil
	}

	pr, err := pullrequest.NewClient(action, gh)
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
			checklists = Parse(action, comment)
		}
	}

	descriptionPR, err := getPullRequestBody(githubContext)
	if err != nil {
		return err
	}

	action.Infof("PR Body: %s", descriptionPR)

	checklists = append(Parse(action, descriptionPR), checklists...)
	return inspect(checklists, action)
}

func inspect(checklists []Checklist, action *githubactions.Action) error {
	action.Debugf("Checklists: %v", checklists)

	action.AddStepSummary("_The following checklists were found and validated:_\n")

	allCompleted := true
	for _, checklist := range checklists {
		allCompleted = allCompleted && checklist.ChecklistCompleted()

		if !checklist.ChecklistCompleted() {
			headerNoPrefix := strings.TrimPrefix(strings.TrimSpace(checklist.Header), "#")
			summaryFormatted := strings.TrimRight(checklist.CompletionSummary(), ".")
			action.Errorf("%s: %s", summaryFormatted, headerNoPrefix)
		}
		action.AddStepSummary(checklist.MarkdownSummary())
	}

	if !allCompleted {
		return errors.New("not all checklists are completed")
	}

	return nil
}

func getPullRequestBody(ghctx *githubactions.GitHubContext) (string, error) {
	pullRequest, ok := ghctx.Event["pull_request"]
	if !ok {
		return "", nil
	}

	body, ok := pullRequest.(map[string]any)["body"]
	if !ok || body == nil {
		return "", nil
	}
	return body.(string), nil
}
