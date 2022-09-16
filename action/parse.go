package action

import (
	"regexp"
	"strings"
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

type reMatch struct {
	LineNumber int
	Raw        string
}

func Parse(content string) Checklist {
	_ = findRE(content, indicatorRE)

	return Checklist{Raw: content}
}

var indicatorRE = regexp.MustCompile(`(?i)<!--\s*Checkmate\s*-->`)
var headerRE = regexp.MustCompile(`(?im)^ {0,3}#{1,6}\s.*`)

func findRE(content string, re *regexp.Regexp) []reMatch {
	var matches []reMatch

	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if re.MatchString(line) {
			matches = append(matches, reMatch{LineNumber: i, Raw: line})
		}
	}

	return matches
}

type block struct {
	Raw         string
	LineNumbers []int
}

func findChecklistBlock(content string) []block {
	re := regexp.MustCompile(`- \[[ x]\] .*`)

	matches := findRE(content, re)

	blocks := []block{}

	b := block{}
	for _, m := range matches {
		// block start
		if b.Raw == "" {
			b.Raw += m.Raw
			b.LineNumbers = append(b.LineNumbers, m.LineNumber)
			continue
		}

		// block continuing
		last := b.LineNumbers[len(b.LineNumbers)-1]
		if last+1 == m.LineNumber {
			b.Raw += "\n" + m.Raw
			b.LineNumbers = append(b.LineNumbers, m.LineNumber)
			continue
		}

		// block ended
		blocks = append(blocks, b)
		b = block{
			Raw:         m.Raw,
			LineNumbers: []int{m.LineNumber},
		}
	}
	if b.Raw != "" {
		blocks = append(blocks, b)
	}

	return blocks
}
