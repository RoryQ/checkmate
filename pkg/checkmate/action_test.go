package checkmate

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/matryer/is"
	"github.com/sethvargo/go-githubactions"
)

func TestRun(t *testing.T) {
	assert := is.NewRelaxed(t)
	require := is.New(t)

	t.Run("CheckedSuccess", func(t *testing.T) {
		action, _ := setupAction("edited-checked")
		err := Run(context.Background(), new(Config), action, nil)
		assert.NoErr(err)
	})

	t.Run("UncheckedFailure", func(t *testing.T) {
		action, _ := setupAction("edited")
		err := Run(context.Background(), new(Config), action, nil)
		require.True(err != nil)
		assert.Equal("not all checklists are completed", err.Error())
	})

	t.Run("OpenedWithNullBody", func(t *testing.T) {
		action, _ := setupAction("opened.with-null-body")
		err := Run(context.Background(), new(Config), action, nil)
		assert.NoErr(err)
	})
}

func setupAction(input string) (*githubactions.Action, *bytes.Buffer) {
	envMap := map[string]string{
		"GITHUB_EVENT_PATH":   fmt.Sprintf("../../test/events/pull-request.%s.json", input),
		"GITHUB_STEP_SUMMARY": "/dev/null",
		"GITHUB_REPOSITORY":   "RoryQ/checkmate",
	}
	getenv := func(key string) string {
		return envMap[key]
	}

	b := new(bytes.Buffer)

	action := githubactions.New(
		githubactions.WithGetenv(getenv),
		githubactions.WithWriter(b),
	)
	return action, b
}
