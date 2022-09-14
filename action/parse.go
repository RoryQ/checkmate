package action

import "regexp"

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
	Start int
	End   int
	Raw   string
}

func Parse(content string) Checklist {

	_ = findRE(content, indicatorRE)

	return Checklist{Raw: content}
}

var indicatorRE = regexp.MustCompile(`(?i)<!--\s*Checkmate\s*-->`)
var headerRE = regexp.MustCompile(`(?im)^ {0,3}#{1,6}\s.*`)

func findRE(content string, re *regexp.Regexp) []reMatch {
	var indexes []reMatch

	for _, match := range re.FindAllStringIndex(content, -1) {
		indexes = append(indexes, reMatch{
			Start: match[0],
			End:   match[1],
			Raw:   content[match[0]:match[1]],
		},
		)

	}

	return indexes
}
