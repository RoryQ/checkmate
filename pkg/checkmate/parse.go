package checkmate

import (
	"regexp"
	"sort"
	"strings"
)

type reMatch struct {
	LineNumber int
	Raw        string
}

var (
	headerRE = regexp.MustCompile(`(?im)^ {0,3}#{1,6}\s.*`)
)

func Parse(content string) (list []Checklist) {
	indicators := findRE(content, indicatorRE)
	if len(indicators) == 0 {
		return
	}

	checklists := findChecklistBlocks(content)

	headers := findRE(content, headerRE)

	for _, ind := range indicators {
		indLineNumber := ind.LineNumber

		c := sort.Search(len(checklists), func(i int) bool {
			return checklists[i].LineNumbers[0] > indLineNumber
		})

		list = append(list, Checklist{
			Items:  blockToItems(checklists[c]),
			Header: closestHeaderTo(headers, indLineNumber),
			Raw:    checklists[c].Raw,
			Meta:   ParseIndicator(ind.Raw),
		})
	}

	return
}

func closestHeaderTo(headers []reMatch, indLineNumber int) string {
	if len(headers) == 0 {
		return ""
	}

	h := sort.Search(len(headers), func(i int) bool {
		return headers[i].LineNumber > indLineNumber
	})

	foundHeader := headers[h-1].Raw
	return foundHeader
}

func blockToItems(b block) (items []ChecklistItem) {
	re := regexp.MustCompile(`- (?P<Checked>\[[ x]]) (?P<Message>.*)`)
	parseChecked := func(s string) bool {
		return s == "[x]"
	}

	for _, line := range strings.Split(b.Raw, "\n") {
		matches := re.FindAllStringSubmatch(line, -1)[0]
		items = append(items,
			ChecklistItem{
				Message: matches[re.SubexpIndex("Message")],
				Checked: parseChecked(matches[re.SubexpIndex("Checked")]),
				Raw:     line,
			},
		)
	}

	return
}

func findRE(content string, re *regexp.Regexp) []reMatch {
	var matches []reMatch

	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if re.MatchString(line) {
			matches = append(matches, reMatch{LineNumber: i, Raw: line})
		}
	}

	sort.SliceStable(matches, func(i, j int) bool { return j > i })
	return matches
}

type block struct {
	Raw         string
	LineNumbers []int
}

func findChecklistBlocks(content string) (blocks []block) {
	re := regexp.MustCompile(`- \[[ x]\] .*`)

	matches := findRE(content, re)

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

var indicatorRE = regexp.MustCompile(`(?i)<!--\s*Checkmate\s*(select=(?P<select>\d+?))?\s*(filepath=(?P<filepath>.+?))?\s*-->`)

func indicatorGroupMatch(s, name string) string {
	match := indicatorRE.FindStringSubmatch(s)
	if i := indicatorRE.SubexpIndex(name); i > 0 && i < len(match) {
		return match[i]
	}
	return ""
}
