package checkmate

import (
	"context"
	"testing"

	"github.com/google/go-github/v48/github"
	"github.com/matryer/is"
	"github.com/migueleliasweb/go-github-mock/src/mock"
)

func Test_commenter(t *testing.T) {
	assert := is.NewRelaxed(t)
	ctx := context.Background()

	const schemaMigrationsGlob = "schema/migrations/*.sql"
	t.Run("NoMatchingFiles", func(t *testing.T) {
		action, _ := setupAction("edited")
		cfg := Config{
			PathsChecklists: map[string]ChecklistsForPath{
				schemaMigrationsGlob: []string{
					"There are no breaking changes in these migrations",
					"I have notified X team of the new schema changes",
				},
			},
		}

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
		err := commenter(ctx, cfg, action, gh)
		assert.NoErr(err)
	})

	t.Run("MatchingFilesNoExistingComment", func(t *testing.T) {
		action, _ := setupAction("edited")
		cfg := Config{
			PathsChecklists: map[string]ChecklistsForPath{
				schemaMigrationsGlob: []string{
					"There are no breaking changes in these migrations",
					"I have notified X team of the new schema changes",
				},
			},
		}

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
				[]github.IssueComment{},
			),
			mock.WithRequestMatch(
				mock.PostReposIssuesCommentsByOwnerByRepoByIssueNumber,
				github.IssueComment{
					Body: github.String(cfg.PathsChecklists[schemaMigrationsGlob].ToChecklistItemsMD(schemaMigrationsGlob)),
				},
			),
		)

		gh := github.NewClient(ghMockAPI)
		err := commenter(ctx, cfg, action, gh)
		assert.NoErr(err)
	})

	t.Run("MatchingFilesWithExistingComment", func(t *testing.T) {
		action, _ := setupAction("edited")
		cfg := Config{
			PathsChecklists: map[string]ChecklistsForPath{
				schemaMigrationsGlob: []string{
					"There are no breaking changes in these migrations",
					"I have notified X team of the new schema changes",
				},
			},
		}

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
						Body: github.String(cfg.PathsChecklists[schemaMigrationsGlob].ToChecklistItemsMD(schemaMigrationsGlob)),
					},
				},
			),
		)

		gh := github.NewClient(ghMockAPI)
		err := commenter(ctx, cfg, action, gh)
		assert.NoErr(err)
	})
}
