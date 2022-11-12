package checkmate

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
	"github.com/sethvargo/go-githubactions"
	"gopkg.in/yaml.v3"
)

type Config struct {
	PathsChecklists map[string]ChecklistsForPath
}

type ChecklistsForPath []string

func (c ChecklistsForPath) ToChecklistItemsMD(glob string) string {
	indicator := fmt.Sprintf("<!-- Checkmate filepath=%s -->\n", glob)
	items := lo.Map(c, func(item string, _ int) string { return "- [ ] " + item })
	return indicator + strings.Join(items, "\n")
}

func ConfigFromInputs(action *githubactions.Action) (*Config, error) {
	c := Config{
		PathsChecklists: map[string]ChecklistsForPath{},
	}
	checklistPaths := action.GetInput("paths")
	if checklistPaths == "" {
		return &c, nil
	}

	if err := yaml.Unmarshal([]byte(checklistPaths), &c.PathsChecklists); err != nil {
		return nil, err
	}

	return &c, nil
}
