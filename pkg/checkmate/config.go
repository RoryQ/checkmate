package checkmate

import (
	"fmt"
	"strings"

	"github.com/niemeyer/pretty"
	"github.com/samber/lo"
	"github.com/sethvargo/go-githubactions"
	"gopkg.in/yaml.v3"

	"github.com/roryq/checkmate/pkg/checkmate/inputs"
)

type Config struct {
	PathsChecklists map[string]ChecklistsForPath
}

type ChecklistsForPath []string

func (c ChecklistsForPath) ToChecklistItemsMD(filenameGlob string) string {
	header := fmt.Sprintf("### Checklist for files matching *%s*\n", strings.ReplaceAll(filenameGlob, "*", "\\*"))
	indicator := fmt.Sprintf("<!-- Checkmate filepath=%s -->\n", filenameGlob)
	items := lo.Map(c, func(item string, _ int) string { return "- [ ] " + item })
	return header + indicator + strings.Join(items, "\n")
}

func ConfigFromInputs(action *githubactions.Action) (*Config, error) {
	action.Infof("Reading Config From Inputs")
	c := Config{
		PathsChecklists: map[string]ChecklistsForPath{},
	}
	checklistPaths := action.GetInput(inputs.Paths)
	if checklistPaths == "" {
		return &c, nil
	}

	if err := yaml.Unmarshal([]byte(checklistPaths), &c.PathsChecklists); err != nil {
		return nil, err
	}

	action.Infof("Config: %s", pretty.Sprint(c))

	return &c, nil
}
