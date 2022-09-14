package action

import (
	"github.com/matryer/is"
	"regexp"
	"strings"
	"testing"
)

var (
	// sample template taken from backstage/backstage
	backstageNoIndicator = `
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

	backstageWithIndicator = insertLine(backstageNoIndicator, "<!-- Checkmate -->", 8)
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
			name: "No Attribute",
			args: args{content: backstageNoIndicator},
			expected: Checklist{
				Raw: backstageNoIndicator,
			},
		},
		{
			name: "With Attribute",
			args: args{content: backstageWithIndicator},
			expected: Checklist{
				Raw:    backstageWithIndicator,
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
						Message: "All your commits have a `Signed-off-by` line in the message. ([more info](https://github.com/backstage/backstage/blob/master/CONTRIBUTING.md#developer-certificate-of-origin))}",
						Checked: false,
						Raw:     "- [ ] All your commits have a `Signed-off-by` line in the message. ([more info](https://github.com/backstage/backstage/blob/master/CONTRIBUTING.md#developer-certificate-of-origin))}",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := is.NewRelaxed(t)
			actual := Parse(tt.args.content)
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
				content: backstageNoIndicator,
				re:      indicatorRE,
			},
			expected: nil,
		},
		{
			name: "WithIndicator",
			args: args{
				content: backstageWithIndicator,
				re:      indicatorRE,
			},
			expected: []reMatch{
				{227, 245, "<!-- Checkmate -->"},
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
				content: backstageWithIndicator,
				re:      headerRE,
			},
			expected: []reMatch{
				{1, 36, "## Hey, I just made a Pull Request!"},
				{192, 225, "#### :heavy_check_mark: Checklist"},
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
