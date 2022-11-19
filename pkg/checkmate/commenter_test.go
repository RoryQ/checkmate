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
	require := is.New(t)
	ctx := context.Background()

	const schemaMigrationsGlob = "schema/migrations/*.sql"
	const docsGlob = "docs/**/*.md"
	const assetsGlob = "assets/**/*.png"
	const selectGlob = "select/**/*.go"
	cfg := Config{
		PathsChecklists: map[string]ChecklistsForPath{
			schemaMigrationsGlob: []string{
				"There are no breaking changes in these migrations",
				"I have notified X team of the new schema changes",
			},
			docsGlob: []string{
				"I have got feedback from the grammar police",
			},
			assetsGlob: []string{
				"Images have been compressed",
			},
			selectGlob: []string{
				"<!--Checkmate select=1-->",
				"Item 1",
				"Item 2",
			},
		},
		Preamble: "Good job!",
	}
	schemaMigrationsChecklist := cfg.PathsChecklists[schemaMigrationsGlob].ToChecklistItemsMD(schemaMigrationsGlob)
	docsChecklist := cfg.PathsChecklists[docsGlob].ToChecklistItemsMD(docsGlob)
	assetsChecklist := cfg.PathsChecklists[assetsGlob].ToChecklistItemsMD(assetsGlob)

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

	t.Run("MatchingFilesForSelectList", func(t *testing.T) {
		action, _ := setupAction("edited")
		ghMockAPI := mock.NewMockedHTTPClient(
			mock.WithRequestMatch(
				mock.GetReposPullsFilesByOwnerByRepoByPullNumber,
				[]github.CommitFile{
					{Filename: ptr.To("select/file/example.go")},
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

					selectList := cfg.PathsChecklists[selectGlob].ToChecklistItemsMD(selectGlob)
					assert.True(strings.Contains(issue.GetBody(), selectList))
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

	t.Run("MatchingFilesWithExistingCommentAndChangedFiles", func(t *testing.T) {
		action, _ := setupAction("edited")
		schemaChecked := strings.ReplaceAll(schemaMigrationsChecklist, "[ ]", "[x]")
		ghMockAPI := mock.NewMockedHTTPClient(
			mock.WithRequestMatch(
				mock.GetReposPullsFilesByOwnerByRepoByPullNumber,
				[]github.CommitFile{
					{Filename: ptr.To("docs/integrations/github/README.md")},
					{Filename: ptr.To("schema/migrations/001_init.sql")},
				},
			),
			// Existing comment was schema checked and assets
			mock.WithRequestMatch(
				mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
				[]github.IssueComment{
					{
						Body: ptr.To(schemaChecked + "\n" + assetsChecklist),
						User: &github.User{ID: ptr.To(GithubActionsBotID)},
					},
				},
			),
			// Expected update is add docs, keep schema checked, remove assets
			mock.WithRequestMatchHandler(
				mock.PatchReposIssuesCommentsByOwnerByRepoByCommentId,
				http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
					b, err := io.ReadAll(r.Body)

					assert.NoErr(err)
					issue := github.IssueComment{}
					assert.NoErr(json.Unmarshal(b, &issue))

					body := issue.GetBody()
					assert.True(strings.Contains(body, cfg.Preamble))
					assert.True(strings.Contains(body, schemaChecked))
					assert.True(strings.Contains(body, docsChecklist))
					assert.True(!strings.Contains(body, assetsChecklist))
				}),
			),
		)

		gh, err := pullrequest.NewClient(action, github.NewClient(ghMockAPI))
		require.NoErr(err)
		_, err = commenter(ctx, cfg, action, gh)
		require.NoErr(err)
	})
}
