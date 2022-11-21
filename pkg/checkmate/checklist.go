package checkmate

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

const (
	redCross  = "❌ "
	greenTick = "✅ "
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

// MarkdownSummary formats the markdown for this checklist on the job summary page.
func (c Checklist) MarkdownSummary() string {
	return "\n" + join(c.CompletionSummary(), markdownEnquote(c.Header), markdownEnquote(c.Raw), "\n")
}

// CompletionSummary formats the validation result of the checklist.
func (c Checklist) CompletionSummary() string {
	checkedCount := lo.CountBy(c.Items, func(i ChecklistItem) bool {
		return i.Checked
	})

	// radio select
	if c.Meta.SelectCount > 0 {
		if checkedCount > c.Meta.SelectCount {
			return fmt.Sprintf(redCross+"Too many items selected - %d selected but expected %d.", checkedCount, c.Meta.SelectCount)
		} else if checkedCount < c.Meta.SelectCount {
			return fmt.Sprintf(redCross+"Incomplete Selectlist - %d out of %d items selected.", checkedCount, c.Meta.SelectCount)
		}
		return fmt.Sprintf(greenTick+"Selectlist Complete - %d out of %d items selected.", checkedCount, c.Meta.SelectCount)
	}

	// regular checklist
	if checkedCount != len(c.Items) {
		return fmt.Sprintf(redCross+"Incomplete Checklist - %d out of %d items checked.", checkedCount, len(c.Items))
	}

	return fmt.Sprintf(greenTick+"Checklist Complete - %d out of %d items checked.", checkedCount, len(c.Items))
}

func markdownEnquote(s string) string {
	return "> " + strings.ReplaceAll(s, "\n", "\n> ")
}

// ParseIndicator will parse the indicator and extract the metadata within.
func ParseIndicator(s string) ChecklistMetadata {
	var selectCount int
	if selectRaw := indicatorGroupMatch(s, "select"); selectRaw != "" {
		var err error
		selectCount, err = strconv.Atoi(selectRaw)
		if err != nil {
			println("error parsing select count", err.Error())
		}
	}

	return ChecklistMetadata{
		RawIndicator: s,
		FilenameGlob: indicatorGroupMatch(s, "filepath"),
		SelectCount:  selectCount,
	}
}
