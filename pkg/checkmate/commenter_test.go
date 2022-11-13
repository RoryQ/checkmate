package checkmate

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-github/v48/github"
	"github.com/matryer/is"
	"github.com/migueleliasweb/go-github-mock/src/mock"

	"github.com/roryq/checkmate/pkg/ptr"
	"github.com/roryq/checkmate/pkg/pullrequest"
)

func Test_commenter(t *testing.T) {
	assert := is.NewRelaxed(t)
	ctx := context.Background()

	const schemaMigrationsGlob = "schema/migrations/*.sql"
	cfg := Config{
		PathsChecklists: map[string]ChecklistsForPath{
			schemaMigrationsGlob: []string{
				"There are no breaking changes in these migrations",
				"I have notified X team of the new schema changes",
			},
		},
	}
	schemaMigrationsChecklist := cfg.PathsChecklists[schemaMigrationsGlob].ToChecklistItemsMD(schemaMigrationsGlob)
	t.Run("NoMatchingFiles", func(t *testing.T) {
		action, _ := setupAction("edited")
		ghMockAPI := mock.NewMockedHTTPClient(
			mock.WithRequestMatch(
				mock.GetReposPullsFilesByOwnerByRepoByPullNumber,
				[]github.CommitFile{
					{Filename: ptr.To("README.md")},
				},
			),
		)

		pr, err := pullrequest.NewClient(action, github.NewClient(ghMockAPI))
		assert.NoErr(err)
		_, err = commenter(ctx, cfg, action, pr)
		assert.NoErr(err)
	})

	t.Run("MatchingFilesNoExistingComment", func(t *testing.T) {
		action, _ := setupAction("edited")
		ghMockAPI := mock.NewMockedHTTPClient(
			mock.WithRequestMatch(
				mock.GetReposPullsFilesByOwnerByRepoByPullNumber,
				[]github.CommitFile{
					{Filename: ptr.To("README.md")},
					{Filename: ptr.To("schema/migrations/001_init.sql")},
				},
			),
			// No existing comment
			mock.WithRequestMatch(
				mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
				[]github.IssueComment{},
			),

			// assert posts comment
			mock.WithRequestMatchHandler(
				mock.PostReposIssuesCommentsByOwnerByRepoByIssueNumber,
				http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
					b, err := io.ReadAll(r.Body)

					assert.NoErr(err)
					issue := github.IssueComment{}

					assert.NoErr(json.Unmarshal(b, &issue))
					assert.True(strings.Contains(issue.GetBody(), schemaMigrationsChecklist))
				}),
			),
		)

		pr, err := pullrequest.NewClient(action, github.NewClient(ghMockAPI))
		assert.NoErr(err)
		_, err = commenter(ctx, cfg, action, pr)
		assert.NoErr(err)
	})

	t.Run("MatchingFilesWithExistingComment", func(t *testing.T) {
		action, _ := setupAction("edited")
		ghMockAPI := mock.NewMockedHTTPClient(
			mock.WithRequestMatch(
				mock.GetReposPullsFilesByOwnerByRepoByPullNumber,
				[]github.CommitFile{
					{Filename: ptr.To("README.md")},
					{Filename: ptr.To("schema/migrations/001_init.sql")},
				},
			),
			mock.WithRequestMatch(
				mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
				[]github.IssueComment{
					{
						Body: ptr.To(schemaMigrationsChecklist),
						User: &github.User{ID: ptr.To(GithubActionsBotID)},
					},
				},
			),
		)

		gh, err := pullrequest.NewClient(action, github.NewClient(ghMockAPI))
		assert.NoErr(err)
		_, err = commenter(ctx, cfg, action, gh)
		assert.NoErr(err)
	})
}
