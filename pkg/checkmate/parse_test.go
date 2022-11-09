package checkmate

import (
	"regexp"
	"strings"
	"testing"

	"github.com/matryer/is"
)

var (
	// sample template taken from backstage/backstage
	ChecklistNoIndicator = `
## Hey, I just made a Pull Request!

<!-- Please describe what you added, and add a screenshot if possible.
     That makes it easier to understand the change so we can :shipit: faster. -->

#### :heavy_check_mark: Checklist

<!--- Please include the following in your Pull Request when applicable: -->

- [ ] A changeset describing the change and affected packages. ([more info](https://github.com/backstage/backstage/blob/master/CONTRIBUTING.md#creating-changesets))
- [ ] Added or updated documentation
- [x] Tests for new functionality and regression tests for bug fixes
- [ ] Screenshots attached (for UI changes)
` + "- [ ] All your commits have a `Signed-off-by` line in the message. ([more info](https://github.com/backstage/backstage/blob/master/CONTRIBUTING.md#developer-certificate-of-origin))"

	Indicator              = "<!-- Checkmate -->"
	ChecklistWithIndicator = insertLine(ChecklistNoIndicator, Indicator, 8)

	checklistOnly = `- [ ] A changeset describing the change and affected packages. ([more info](https://github.com/backstage/backstage/blob/master/CONTRIBUTING.md#creating-changesets))
- [ ] Added or updated documentation
- [x] Tests for new functionality and regression tests for bug fixes
- [ ] Screenshots attached (for UI changes)
` + "- [ ] All your commits have a `Signed-off-by` line in the message. ([more info](https://github.com/backstage/backstage/blob/master/CONTRIBUTING.md#developer-certificate-of-origin))"
)

func TestParse(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name     string
		args     args
		expected Checklist
	}{
		{
			name:     "Empty",
			args:     args{},
			expected: Checklist{},
		},
		{
			name: "No Indicator",
			args: args{content: ChecklistNoIndicator},
			expected: Checklist{
				Raw: ChecklistNoIndicator,
			},
		},
		{
			name: "With Indicator",
			args: args{content: ChecklistWithIndicator},
			expected: Checklist{
				Raw:    checklistOnly,
				Header: "#### :heavy_check_mark: Checklist",
				Items: []ChecklistItem{
					{
						Message: "A changeset describing the change and affected packages. ([more info](https://github.com/backstage/backstage/blob/master/CONTRIBUTING.md#creating-changesets))",
						Checked: false,
						Raw:     `- [ ] A changeset describing the change and affected packages. ([more info](https://github.com/backstage/backstage/blob/master/CONTRIBUTING.md#creating-changesets))`,
					},
					{
						Message: "Added or updated documentation",
						Checked: false,
						Raw:     `- [ ] Added or updated documentation`,
					},
					{
						Message: "Tests for new functionality and regression tests for bug fixes",
						Checked: true,
						Raw:     `- [x] Tests for new functionality and regression tests for bug fixes`,
					},
					{
						Message: "Screenshots attached (for UI changes)",
						Checked: false,
						Raw:     `- [ ] Screenshots attached (for UI changes)`,
					},
					{
						Message: "All your commits have a `Signed-off-by` line in the message. ([more info](https://github.com/backstage/backstage/blob/master/CONTRIBUTING.md#developer-certificate-of-origin))",
						Checked: false,
						Raw:     "- [ ] All your commits have a `Signed-off-by` line in the message. ([more info](https://github.com/backstage/backstage/blob/master/CONTRIBUTING.md#developer-certificate-of-origin))",
					},
				},
			},
		},
		{
			name: "No Headers",
			args: args{content: "<!--Checkmate-->\n- [ ] unchecked"},
			expected: Checklist{
				Header: "",
				Raw:    "- [ ] unchecked",
				Items: []ChecklistItem{
					{
						"unchecked",
						false,
						"- [ ] unchecked",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := is.NewRelaxed(t)
			actualList := Parse(tt.args.content)
			for _, actual := range actualList {
				assert.Equal(actual.Raw, tt.expected.Raw)
				assert.Equal(actual.Header, tt.expected.Header)

				for i := range tt.expected.Items {
					var actualItem ChecklistItem
					if len(actual.Items) > i {
						actualItem = actual.Items[i]
					}
					expectedItem := tt.expected.Items[i]
					assert.Equal(actualItem.Message, expectedItem.Message)
					assert.Equal(actualItem.Checked, expectedItem.Checked)
					assert.Equal(actualItem.Raw, expectedItem.Raw)
				}
			}
		})
	}
}

func insertLine(all, insert string, index int) string {
	toLines := strings.Split(all, "\n")
	toInsert := []string{insert}
	concatAll := append(toLines[:index], append(toInsert, toLines[index:]...)...)
	return strings.Join(concatAll, "\n")
}

func TestRegexps(t *testing.T) {
	type args struct {
		content string
		re      *regexp.Regexp
	}
	tests := []struct {
		name     string
		args     args
		expected []reMatch
	}{
		{
			name: "NoIndicator",
			args: args{
				content: ChecklistNoIndicator,
				re:      indicatorRE,
			},
			expected: nil,
		},
		{
			name: "WithIndicator",
			args: args{
				content: ChecklistWithIndicator,
				re:      indicatorRE,
			},
			expected: []reMatch{
				{8, "<!-- Checkmate -->"},
			},
		},
		{
			name: "NoHeader",
			args: args{
				content: "",
				re:      headerRE,
			},
			expected: nil,
		},
		{
			name: "WithHeader",
			args: args{
				content: ChecklistWithIndicator,
				re:      headerRE,
			},
			expected: []reMatch{
				{1, "## Hey, I just made a Pull Request!"},
				{6, "#### :heavy_check_mark: Checklist"},
			},
		},
		{
			name: "IndentedCodeNotHeader",
			args: args{
				content: "    # code indented text",
				re:      headerRE,
			},
			expected: nil,
		},
		{
			name: "TextNotHeader",
			args: args{
				content: "Text with a ## inside",
				re:      headerRE,
			},
			expected: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := is.NewRelaxed(t)
			actual := findRE(tt.args.content, tt.args.re)
			assert.Equal(tt.expected, actual)
		})
	}
}

func Test_findChecklistBlock(t *testing.T) {
	type args struct {
		content string
	}
	checklistBlock := `- [ ] A changeset describing the change and affected packages. ([more info](https://github.com/backstage/backstage/blob/master/CONTRIBUTING.md#creating-changesets))
- [ ] Added or updated documentation
- [x] Tests for new functionality and regression tests for bug fixes
- [ ] Screenshots attached (for UI changes)
` + "- [ ] All your commits have a `Signed-off-by` line in the message. ([more info](https://github.com/backstage/backstage/blob/master/CONTRIBUTING.md#developer-certificate-of-origin))"
	tests := []struct {
		name     string
		args     args
		expected []block
	}{
		{
			name:     "Empty",
			args:     args{},
			expected: []block{},
		},
		{
			name: "OneBlock",
			args: args{
				content: ChecklistWithIndicator,
			},
			expected: []block{
				{
					Raw:         checklistBlock,
					LineNumbers: []int{11, 12, 13, 14, 15},
				},
			},
		},
		{
			name: "MultiBlock",
			args: args{
				content: ChecklistWithIndicator + ChecklistNoIndicator,
			},
			expected: []block{
				{
					Raw:         checklistBlock,
					LineNumbers: []int{11, 12, 13, 14, 15},
				},
				{
					Raw:         checklistBlock,
					LineNumbers: []int{25, 26, 27, 28, 29},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := is.NewRelaxed(t)
			actual := findChecklistBlocks(tt.args.content)

			for i, expectedB := range tt.expected {
				var actualB block
				if len(actual) > i {
					actualB = actual[i]
				}

				assert.Equal(actualB.Raw, expectedB.Raw)
				assert.Equal(actualB.LineNumbers, expectedB.LineNumbers)
			}
		})
	}
}
