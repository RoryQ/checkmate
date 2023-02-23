package pullrequest

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/sethvargo/go-githubactions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_getPRNumber(t *testing.T) {
	readEvent := func(eventName string) map[string]any {
		eventData, err := os.ReadFile(fmt.Sprintf("../../test/events/%s.json", eventName))
		require.NoError(t, err)
		var githubContext githubactions.GitHubContext
		require.NoError(t, json.Unmarshal(eventData, &githubContext.Event))
		return githubContext.Event
	}

	t.Run("pull_request_event", func(t *testing.T) {
		prNum, err := getPRNumber(readEvent("pull-request.edited"))
		assert.NoError(t, err)
		assert.Equal(t, 2, prNum)
	})

	t.Run("issue_comment_event", func(t *testing.T) {
		prNum, err := getPRNumber(readEvent("issue-comment.created"))
		assert.NoError(t, err)
		assert.Equal(t, 1, prNum)
	})
}
