package checkmate

import (
	"testing"

	"github.com/google/go-cmp/cmp"
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
