package checkmate

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/google/go-github/v48/github"
	"github.com/matryer/is"
	"github.com/migueleliasweb/go-github-mock/src/mock"
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
					{Filename: github.String("README.md")},
					//{Filename: github.String("schema/migrations/001_init.sql")},
				},
			),
		)

		gh := github.NewClient(ghMockAPI)
		_, err := commenter(ctx, cfg, action, gh)
		assert.NoErr(err)
	})

	t.Run("MatchingFilesNoExistingComment", func(t *testing.T) {
		action, _ := setupAction("edited")
		ghMockAPI := mock.NewMockedHTTPClient(
			mock.WithRequestMatch(
				mock.GetReposPullsFilesByOwnerByRepoByPullNumber,
				[]github.CommitFile{
					{Filename: github.String("README.md")},
					{Filename: github.String("schema/migrations/001_init.sql")},
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
					assert.Equal(*issue.Body, schemaMigrationsChecklist)
				}),
			),
		)

		gh := github.NewClient(ghMockAPI)
		_, err := commenter(ctx, cfg, action, gh)
		assert.NoErr(err)
	})

	t.Run("MatchingFilesWithExistingComment", func(t *testing.T) {
		action, _ := setupAction("edited")
		ghMockAPI := mock.NewMockedHTTPClient(
			mock.WithRequestMatch(
				mock.GetReposPullsFilesByOwnerByRepoByPullNumber,
				[]github.CommitFile{
					{Filename: github.String("README.md")},
					{Filename: github.String("schema/migrations/001_init.sql")},
				},
			),
			mock.WithRequestMatch(
				mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
				[]github.IssueComment{
					{
						Body: github.String(schemaMigrationsChecklist),
					},
				},
			),
		)

		gh := github.NewClient(ghMockAPI)
		_, err := commenter(ctx, cfg, action, gh)
		assert.NoErr(err)
	})
}
