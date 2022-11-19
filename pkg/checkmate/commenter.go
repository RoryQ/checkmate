package checkmate

import (
	"context"
	"regexp"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v48/github"
	"github.com/samber/lo"
	"github.com/sethvargo/go-githubactions"

	"github.com/roryq/checkmate/pkg/ptr"
	"github.com/roryq/checkmate/pkg/pullrequest"
)

func commenter(ctx context.Context, cfg Config, action *githubactions.Action, pr pullrequest.Client) (string, error) {
	fileNames, err := listPullRequestFiles(ctx, pr)
	if err != nil {
		return "", err
	}

	action.Debugf("checking files: \n%s", strings.Join(fileNames, "\n"))
	matched := lo.Filter(lo.Keys(cfg.PathsChecklists), func(pathGlob string, _ int) bool {
		for _, name := range fileNames {
			if matched, _ := doublestar.Match(pathGlob, name); matched {
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

// updateComment will
// 1. Create a new checklist comment when there is no existing comment
// 2. Update the existing checklist based on the configured filename globs. Checklists will be removed if there is not
// a matching glob, new blank checklists will be added for new matches and existing checklists will remain as is.
func updateComment(ctx context.Context, action *githubactions.Action, cfg Config, pr pullrequest.Client, checklistConfig map[string]ChecklistsForPath, comment *github.IssueComment) (string, error) {
	keys := sorted(lo.Keys(checklistConfig))

	// No existing comment
	if comment.GetBody() == "" {
		action.Infof("Writing new automated checklist")

		allChecklists := strings.Join(lo.Map(keys, func(k string, _ int) string {
			return checklistConfig[k].ToChecklistItemsMD(k)
		}), "\n\n")

		commentBody := join(cfg.Preamble, allChecklists, "\n")

		_, err := pr.CreateComment(ctx, &github.IssueComment{Body: ptr.To(commentBody)})
		return commentBody, err
	}

	checklists := Parse(comment.GetBody())
	actualFilenames := sorted(lo.WithoutEmpty(lo.Map(checklists, func(item Checklist, _ int) string {
		return item.Meta.FilenameGlob
	})))

	// Existing comment has no filename glob changes
	if cmp.Equal(keys, actualFilenames) {
		action.Infof("No changes with existing comment")
		return comment.GetBody(), nil
	}
	action.Infof("Changes detected for existing comment")

	globAdd, globRemove := lo.Difference(keys, actualFilenames)
	action.Infof("Removing checklists for [ %s ]", strings.Join(globRemove, " "))

	filtered := lo.Filter(checklists, func(item Checklist, _ int) bool {
		return !lo.Contains(globRemove, item.Meta.FilenameGlob)
	})

	action.Infof("Adding checklists for [ %s ]", strings.Join(globAdd, " "))
	toAdd := lo.Map(globAdd, func(key string, _ int) Checklist {
		return Parse(checklistConfig[key].ToChecklistItemsMD(key))[0]
	})

	checklists = sortedByFilename(append(filtered, toAdd...))

	allChecklists := strings.Join(lo.Map(checklists, func(it Checklist, _ int) string {
		return join(it.Header, it.Meta.RawIndicator, it.Raw, "\n")
	}), "\n\n")

	editedComment := github.IssueComment{
		ID:   comment.ID,
		Body: ptr.To(join(cfg.Preamble, allChecklists, "\n")),
	}

	action.Infof("Editing existing comment")
	// Add new checklists and remove checklists for missing fileglob matches
	if _, err := pr.EditComment(ctx, &editedComment); err != nil {
		return "", err
	}

	return editedComment.GetBody(), nil
}

func listPullRequestFiles(ctx context.Context, pr pullrequest.Client) ([]string, error) {
	files, err := pr.ListFiles(ctx, nil)
	if err != nil {
		return nil, err
	}
	fileNames := lo.Map(files, func(item *github.CommitFile, _ int) string { return item.GetFilename() })
	return fileNames, nil
}

func sortedByFilename(cs []Checklist) []Checklist {
	sort.SliceStable(cs, func(i, j int) bool {
		return cs[i].Meta.FilenameGlob < cs[j].Meta.FilenameGlob
	})
	return cs
}

func sorted(ss []string) []string {
	sort.StringSlice(ss).Sort()
	return ss
}

// join call strings.Join using the last argument as the separator
func join(s ...string) string {
	return strings.Join(s[:len(s)-1], s[len(s)-1])
}
