package checkmate

import (
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sethvargo/go-githubactions"
	"github.com/stretchr/testify/assert"
)

func TestParseIndicator(t *testing.T) {
	tests := []struct {
		Name      string
		Indicator string
		Want      ChecklistMetadata
	}{
		// TODO: Add test cases.
		{
			Name:      "WithFilepathGlob",
			Indicator: "<!-- Checkmate filepath=schema/migrations/*.sql -->",
			Want: ChecklistMetadata{
				RawIndicator: "<!-- Checkmate filepath=schema/migrations/*.sql -->",
				FilenameGlob: "schema/migrations/*.sql",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if got := ParseIndicator(tt.Indicator); !cmp.Equal(got, tt.Want) {
				t.Error(cmp.Diff(got, tt.Want))
			}
		})
	}
}

func TestChecklist_MarkdownSummary(t *testing.T) {
	tests := []struct {
		name         string
		checklistRaw string
		want         string
	}{
		{
			name: "Selectlist Complete",
			checklistRaw: `
## Choose One Item
<!--Checkmate select=1-->
- [ ] Apple
- [x] Pear`,
			want: `
✅ Selectlist Complete - 1 out of 1 items selected.
> ## Choose One Item
> - [ ] Apple
> - [x] Pear`,
		},
		{
			name: "Selectlist Invalid",
			checklistRaw: `
## Choose One Item
<!--Checkmate select=1-->
- [x] Apple
- [x] Pear`,
			want: `
❌ Too many items selected - 2 selected but expected 1.
> ## Choose One Item
> - [x] Apple
> - [x] Pear`,
		},
		{
			name: "Selectlist Incomplete",
			checklistRaw: `
## Choose One Item
<!--Checkmate select=1-->
- [ ] Apple
- [ ] Pear`,
			want: `
❌ Incomplete Selectlist - 0 out of 1 items selected.
> ## Choose One Item
> - [ ] Apple
> - [ ] Pear`,
		},
		{
			name: "Checklist Complete",
			checklistRaw: `
## My Checklist
<!--Checkmate-->
- [x] Apple
- [x] Pear`,
			want: `
✅ Checklist Complete - 2 out of 2 items checked.
> ## My Checklist
> - [x] Apple
> - [x] Pear`,
		},
		{
			name: "Checklist Incomplete",
			checklistRaw: `
## My Checklist
<!--Checkmate-->
- [x] Apple
- [ ] Pear`,
			want: `
❌ Incomplete Checklist - 1 out of 2 items checked.
> ## My Checklist
> - [x] Apple
> - [ ] Pear`,
		},
	}
	action := githubactions.New(githubactions.WithWriter(io.Discard))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checklists := Parse(action, tt.checklistRaw)
			assert.Equal(t, len(checklists), 1)
			c := checklists[0]
			if got := c.MarkdownSummary(); got != tt.want {
				t.Errorf("MarkdownSummary() = %v, want %v", got, tt.want)
			}
		})
	}
}
