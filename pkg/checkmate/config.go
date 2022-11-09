package checkmate

import (
	"encoding/json"

	"github.com/sethvargo/go-githubactions"
)

type Config struct {
	PathsChecklists map[string]string
}

func ConfigFromInputs(action *githubactions.Action) (*Config, error) {
	c := Config{}
	checklistsJson := action.GetInput("paths")
	if checklistsJson == "" {
		return &c, nil
	}

	checklists := make(map[string]string)
	err := json.Unmarshal([]byte(checklistsJson), &checklists)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
