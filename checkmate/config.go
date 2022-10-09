package checkmate

import (
	"encoding/json"
	"log"

	githubactions "github.com/sethvargo/go-githubactions"
)

type Config struct {
}

func ConfigFromInputs(action *githubactions.Action) (*Config, error) {
	checklistsJson := action.GetInput("checklists")

	log.Print(checklistsJson)
	checklists := make(map[string]string)
	err := json.Unmarshal([]byte(checklistsJson), &checklists)
	if err != nil {
		return nil, err
	}

	c := Config{}
	return &c, nil
}
