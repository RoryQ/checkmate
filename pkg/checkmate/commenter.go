package checkmate

import (
	"context"
	"errors"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/google/go-github/v48/github"
	"github.com/samber/lo"
	"github.com/sethvargo/go-githubactions"
)

func commenter(ctx context.Context, cfg Config, action *githubactions.Action, gh *github.Client) (string, error) {
	pr, err := getPullRequestContext(action)
	if err != nil {
		return "", err
	}

	fileNames, err := listPullRequestFiles(ctx, gh, pr)
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

	comment, err := getExistingComment(ctx, gh, pr)
	if err != nil {
		return "", err
	}

	return updateComment(ctx, action, gh, pr, checklists, comment)
}

var commenterIndicatorRE = regexp.MustCompile(`(?i)<!--\s*Checkmate\s+filepath=.*?-->`)

func getExistingComment(ctx context.Context, gh *github.Client, pr pullRequestContext) (*github.IssueComment, error) {
	comments, _, err := gh.Issues.ListComments(ctx, pr.Owner, pr.Repo, pr.Number, nil)
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

func isBotID(id int64) bool {
	const GithubActionsBotID = 41898282
	return id == GithubActionsBotID
}

func updateComment(ctx context.Context, action *githubactions.Action, gh *github.Client, pr pullRequestContext, checklists map[string]ChecklistsForPath, comment *github.IssueComment) (string, error) {
	keys := lo.Keys(checklists)
	sort.StringSlice(keys).Sort()

	if comment.GetBody() == "" {
		action.Infof("Writing new automated checklist")

		allChecklists := strings.Join(lo.Map(keys, func(k string, _ int) string {
			return checklists[k].ToChecklistItemsMD(k)
		}), "\n\n")

		preamble := "Thanks for your contribution! Please complete the following tasks related to your changes and tick" +
			"the checklists when complete."

		commentBody := preamble + "\n" + allChecklists
		_, _, err := gh.Issues.CreateComment(ctx, pr.Owner, pr.Repo, pr.Number, &github.IssueComment{
			Body: github.String(commentBody),
		})
		return commentBody, err
	}

	// TODO add / remove checklists based on file changes
	return comment.GetBody(), nil
}

type pullRequestContext struct {
	Owner  string
	Repo   string
	Number int
}

func getPullRequestContext(action *githubactions.Action) (pullRequestContext, error) {
	ctx, err := action.Context()
	if err != nil {
		return pullRequestContext{}, err
	}
	owner, repo := getRepo(action, ctx.Event)
	number, err := getPRNumber(ctx.Event)
	if err != nil {
		return pullRequestContext{}, err
	}

	return pullRequestContext{
		Owner:  owner,
		Repo:   repo,
		Number: number,
	}, nil
}

func listPullRequestFiles(ctx context.Context, gh *github.Client, pr pullRequestContext) ([]string, error) {
	files, _, err := gh.PullRequests.ListFiles(ctx, pr.Owner, pr.Repo, pr.Number, nil)
	if err != nil {
		return nil, err
	}
	fileNames := lo.Map(files, func(item *github.CommitFile, _ int) string { return item.GetFilename() })
	return fileNames, nil
}

func getPRNumber(event map[string]any) (int, error) {
	number, ok := event["pull_request"].(map[string]any)["number"]
	if !ok {
		return 0, errors.New("cannot get pull_request number")
	}
	return int(number.(float64)), nil
}

func getRepo(action *githubactions.Action, event map[string]any) (string, string) {
	splitRepo := func(name string) (string, string) {
		split := strings.Split(name, "/")
		return split[0], split[1]
	}

	if fullName := action.Getenv("GITHUB_REPOSITORY"); fullName != "" {
		splitRepo(fullName)
	}

	if fullName, ok := event["repository"].(map[string]any)["full_name"]; ok {
		return splitRepo(fullName.(string))
	}
	return "", ""
}
