package checkmate

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/samber/lo"
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
	SelectCount  int
}

func (c Checklist) ChecklistCompleted() bool {
	checkedCount := lo.CountBy(c.Items, func(i ChecklistItem) bool {
		return i.Checked
	})

	// if select is enabled then check count
	if c.Meta.SelectCount > 0 {
		return c.Meta.SelectCount == checkedCount
	}

	return len(c.Items) == checkedCount
}

// https://github.blog/2022-05-09-supercharging-github-actions-with-job-summaries/

func (c Checklist) MarkdownSummary() string {
	return join(c.CompletionSummary(), markdownEnquote(c.Header), markdownEnquote(c.Raw), "\n")
}

func (c Checklist) CompletionSummary() string {
	checkedCount := lo.CountBy(c.Items, func(i ChecklistItem) bool {
		return i.Checked
	})

	// radio select
	if c.Meta.SelectCount > 0 {
		if checkedCount > c.Meta.SelectCount {
			return fmt.Sprintf("Too many items selected - %d selected but expected %d.", checkedCount, c.Meta.SelectCount)
		} else if checkedCount < c.Meta.SelectCount {
			return fmt.Sprintf("Incomplete SelectList - %d out of %d items selected.", checkedCount, c.Meta.SelectCount)
		}
		return fmt.Sprintf("SelectList Complete - %d out of %d items selected.", checkedCount, c.Meta.SelectCount)
	}

	// regular checklist
	if checkedCount != len(c.Items) {
		return fmt.Sprintf("Incomplete Checklist - %d out of %d items checked.", checkedCount, len(c.Items))
	}

	return fmt.Sprintf("Checklist Complete - %d out of %d items checked.", checkedCount, len(c.Items))
}

func markdownEnquote(s string) string {
	return "> " + strings.ReplaceAll(s, "\n", "\n> ")
}

func ParseIndicator(s string) ChecklistMetadata {
	match := indicatorRE.FindStringSubmatch(s)

	namedGroupMatch := func(name string) string {
		if i := indicatorRE.SubexpIndex(name); i > 0 {
			return match[i]
		}
		return ""
	}

	var selectCount int
	if selectRaw := namedGroupMatch("select"); selectRaw != "" {
		var err error
		selectCount, err = strconv.Atoi(selectRaw)
		if err != nil {
			println("error parsing select count", err.Error())
		}
	}

	return ChecklistMetadata{
		RawIndicator: s,
		FilenameGlob: namedGroupMatch("filepath"),
		SelectCount:  selectCount,
	}
}
