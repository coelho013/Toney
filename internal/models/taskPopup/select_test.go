package taskpopup

import (
	"strings"
	"testing"

	"github.com/SourcewareLab/Toney/v2/internal/enums"
)

func TestGetText(t *testing.T) {
	titleMap := map[enums.TaskStatus]string{
		enums.Pending:   "Pending",
		enums.Started:   "Started",
		enums.Complete:  "Complete",
		enums.Abandoned: "Abandoned",
	}

	tests := []struct {
		name           string
		opts           []enums.TaskStatus
		selected       int
		expectedLabels []string
	}{
		{
			name:           "default order",
			opts:           []enums.TaskStatus{enums.Pending, enums.Started, enums.Abandoned, enums.Complete},
			selected:       0,
			expectedLabels: []string{"Pending", "Started", "Abandoned", "Complete"},
		},
		{
			name:           "non-default order",
			opts:           []enums.TaskStatus{enums.Complete, enums.Started, enums.Pending, enums.Abandoned},
			selected:       0,
			expectedLabels: []string{"Complete", "Started", "Pending", "Abandoned"},
		},
		{
			name:           "selected last item",
			opts:           []enums.TaskStatus{enums.Pending, enums.Started, enums.Abandoned, enums.Complete},
			selected:       3,
			expectedLabels: []string{"Pending", "Started", "Abandoned", "Complete"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SelectStatus{
				Width:    20,
				Height:   10,
				Opts:     tt.opts,
				TitleMap: titleMap,
				Selected: tt.selected,
			}

			got := s.GetText()
			lines := strings.Split(got, "\n")

			if len(lines) != len(tt.expectedLabels) {
				t.Fatalf("GetText() returned %d lines, want %d", len(lines), len(tt.expectedLabels))
			}

			for i, expected := range tt.expectedLabels {
				if !strings.Contains(lines[i], expected) {
					t.Errorf("GetText() line %d = %q, want it to contain %q", i, lines[i], expected)
				}
			}
		})
	}
}
