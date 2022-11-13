package checkmate

import (
	"context"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/google/go-github/v48/github"
	"github.com/samber/lo"
	"github.com/sethvargo/go-githubactions"

	"github.com/roryq/checkmate/pkg/pullrequest"
)

func commenter(ctx context.Context, cfg Config, action *githubactions.Action, pr pullrequest.Client) (string, error) {
	fileNames, err := listPullRequestFiles(ctx, pr)
	if err != nil {
		return "", err
	}

	matched := lo.Filter(lo.Keys(cfg.PathsChecklists), func(pathGlob string, _ int) bool {
		for _, name := range fileNames {
			if matched, _ := filepath.Match(pathGlob, name); matched {
				return true
			}
		}
		return false
	})

	if len(matched) == 0 {
		action.Infof("no matched paths")
		return "", nil
	}

	action.Infof("matched paths: [ %s ]", strings.Join(matched, " "))
	checklists := lo.PickByKeys(cfg.PathsChecklists, matched)

	comment, err := getExistingComment(ctx, pr)
	if err != nil {
		return "", err
	}

	return updateComment(ctx, action, cfg, pr, checklists, comment)
}

var commenterIndicatorRE = regexp.MustCompile(`(?i)<!--\s*Checkmate\s+filepath=.*?-->`)

func getExistingComment(ctx context.Context, pr pullrequest.Client) (*github.IssueComment, error) {
	comments, err := pr.ListComments(ctx, nil)
	if err != nil {
		return nil, err
	}

	comments = lo.Filter(comments, func(c *github.IssueComment, _ int) bool {
		return isBotID(c.GetUser().GetID()) &&
			commenterIndicatorRE.MatchString(c.GetBody())
	})

	if len(comments) == 0 {
		return nil, nil
	}

	return comments[0], nil
}

const GithubActionsBotID int64 = 41898282

func isBotID(id int64) bool {
	return id == GithubActionsBotID
}

func updateComment(ctx context.Context, action *githubactions.Action, cfg Config, pr pullrequest.Client, checklists map[string]ChecklistsForPath, comment *github.IssueComment) (string, error) {
	keys := lo.Keys(checklists)
	sort.StringSlice(keys).Sort()

	if comment.GetBody() == "" {
		action.Infof("Writing new automated checklist")

		allChecklists := strings.Join(lo.Map(keys, func(k string, _ int) string {
			return checklists[k].ToChecklistItemsMD(k)
		}), "\n\n")

		commentBody := cfg.Preamble + "\n" + allChecklists

		_, err := pr.CreateComment(ctx, &github.IssueComment{Body: github.String(commentBody)})
		return commentBody, err
	}

	// TODO add / remove checklists based on file changes
	return comment.GetBody(), nil
}

func listPullRequestFiles(ctx context.Context, pr pullrequest.Client) ([]string, error) {
	files, err := pr.ListFiles(ctx, nil)
	if err != nil {
		return nil, err
	}
	fileNames := lo.Map(files, func(item *github.CommitFile, _ int) string { return item.GetFilename() })
	return fileNames, nil
}
