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
	Meta   ChecklistMetadata
	Header string
	Raw    string
}

type ChecklistMetadata struct {
	RawIndicator string
	FilenameGlob string
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

func ParseIndicator(s string) ChecklistMetadata {
	match := indicatorRE.FindStringSubmatch(s)

	namedGroupMatch := func(name string) string {
		if i := indicatorRE.SubexpIndex(name); i > 0 {
			return match[i]
		}
		return ""
	}

	return ChecklistMetadata{
		RawIndicator: s,
		FilenameGlob: namedGroupMatch("filepath"),
	}
}
