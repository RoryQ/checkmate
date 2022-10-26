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

	t.Run("CheckedSuccess", func(t *testing.T) {
		action := setupAction("edited-checked")
		err := Run(context.Background(), new(Config), action)
		assert.NoErr(err)
	})

	t.Run("UncheckedFailure", func(t *testing.T) {
		action := setupAction("edited")
		err := Run(context.Background(), new(Config), action)
		assert.Equal("not all checklists are completed", err.Error())
	})
}

func setupAction(input string) *githubactions.Action {
	envMap := map[string]string{
		"GITHUB_EVENT_PATH":   fmt.Sprintf("../test/events/pull-request.%s.json", input),
		"GITHUB_STEP_SUMMARY": "/dev/null",
	}
	getenv := func(key string) string {
		return envMap[key]
	}

	b := new(bytes.Buffer)

	action := githubactions.New(
		githubactions.WithGetenv(getenv),
		githubactions.WithWriter(b),
	)
	return action
}
