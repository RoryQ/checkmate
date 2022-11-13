package pullrequest

import (
	"context"
	"errors"
	"strings"

	"github.com/google/go-github/v48/github"
	"github.com/sethvargo/go-githubactions"
)

type Client struct {
	Owner  string
	Repo   string
	Number int
	gh     *github.Client
}

func (pr Client) ListFiles(ctx context.Context, options *github.ListOptions) ([]*github.CommitFile, error) {
	files, _, err := pr.gh.PullRequests.ListFiles(ctx, pr.Owner, pr.Repo, pr.Number, options)
	return files, err
}

func (pr Client) ListComments(ctx context.Context, options *github.IssueListCommentsOptions) ([]*github.IssueComment, error) {
	comments, _, err := pr.gh.Issues.ListComments(ctx, pr.Owner, pr.Repo, pr.Number, options)
	return comments, err
}

func (pr Client) CreateComment(ctx context.Context, comment *github.IssueComment) (*github.IssueComment, error) {
	comment, _, err := pr.gh.Issues.CreateComment(ctx, pr.Owner, pr.Repo, pr.Number, comment)
	return comment, err
}

func NewClient(action *githubactions.Action, gh *github.Client) (Client, error) {
	ctx, err := action.Context()
	if err != nil {
		return Client{}, err
	}
	owner, repo := getRepo(action, ctx.Event)
	number, err := getPRNumber(ctx.Event)
	if err != nil {
		return Client{}, err
	}

	return Client{
		Owner:  owner,
		Repo:   repo,
		Number: number,
		gh:     gh,
	}, nil
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
