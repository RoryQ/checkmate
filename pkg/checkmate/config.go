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
	Preamble        string
}

type ChecklistsForPath []string

func (c ChecklistsForPath) ToChecklistItemsMD(filenameGlob string) string {
	if len(c) == 0 {
		return ""
	}
	header := fmt.Sprintf("### Checklist for files matching *%s*\n", strings.ReplaceAll(filenameGlob, "*", "\\*"))

	var indicator string
	if selectCount := indicatorGroupMatch(c[0], "select"); selectCount != "" {
		indicator = fmt.Sprintf("<!-- Checkmate select=%s filepath=%s -->\n", selectCount, filenameGlob)
		items := lo.Map(c[1:], func(item string, _ int) string { return "- [ ] " + item })
		return header + indicator + strings.Join(items, "\n")
	}

	indicator = fmt.Sprintf("<!-- Checkmate filepath=%s -->\n", filenameGlob)
	items := lo.Map(c, func(item string, _ int) string { return "- [ ] " + item })
	return header + indicator + strings.Join(items, "\n")
}

const PreambleDefaultMessage = "Thanks for your contribution!\n Please complete the following tasks related to your changes and tick " +
	"the checklists when complete."

func ConfigFromInputs(action *githubactions.Action) (*Config, error) {
	action.Infof("Reading Config From Inputs")
	c := Config{
		PathsChecklists: map[string]ChecklistsForPath{},
	}
	checklistPaths := action.GetInput(inputs.Paths)
	if checklistPaths == "" {
		return &c, nil
	}

	c.Preamble = action.GetInput(inputs.Preamble)
	if c.Preamble == "" {
		c.Preamble = PreambleDefaultMessage
	}

	if err := yaml.Unmarshal([]byte(checklistPaths), &c.PathsChecklists); err != nil {
		return nil, err
	}

	action.Infof("Config: %s", pretty.Sprint(c))

	return &c, nil
}
