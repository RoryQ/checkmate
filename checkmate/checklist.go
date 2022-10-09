package checkmate

import (
	"fmt"
)

type ChecklistItem struct {
	Message string
	Checked bool
	Raw     string
}

type Checklist struct {
	Items  []ChecklistItem
	Header string
	Raw    string
}

func (c Checklist) AllChecked() bool {
	for _, item := range c.Items {
		if !item.Checked {
			return false
		}
	}
	return true
}

// https://github.blog/2022-05-09-supercharging-github-actions-with-job-summaries/

func (c Checklist) Summary() string {
	return fmt.Sprintf("%s\n%s", c.Header, c.Raw)
}
