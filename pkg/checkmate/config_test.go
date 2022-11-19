package checkmate

import "testing"

func TestChecklistsForPath_ToChecklistItemsMD(t *testing.T) {
	tests := []struct {
		name         string
		c            ChecklistsForPath
		filenameGlob string
		want         string
	}{
		{
			name:         "NoItems",
			c:            []string{},
			filenameGlob: "**/*.go",
			want:         "",
		},
		{
			name: "RegularChecklist",
			c: []string{
				"item 1",
				"item 2",
			},
			filenameGlob: "**/*.go",
			want: `### Checklist for files matching *\*\*/\*.go*
<!-- Checkmate filepath=**/*.go -->
- [ ] item 1
- [ ] item 2`,
		},
		{
			name: "ContainsSelectListIndicator",
			c: []string{
				"<!-- checkmate select=1 -->",
				"item 1",
				"item 2",
			},
			filenameGlob: "**/*.go",
			want: `### Checklist for files matching *\*\*/\*.go*
<!-- Checkmate select=1 filepath=**/*.go -->
- [ ] item 1
- [ ] item 2`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.ToChecklistItemsMD(tt.filenameGlob); got != tt.want {
				t.Errorf("ToChecklistItemsMD() = %v, want %v", got, tt.want)
			}
		})
	}
}
