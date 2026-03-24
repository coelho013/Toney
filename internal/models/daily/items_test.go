package daily

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/SourcewareLab/Toney/v2/internal/enums"
)

func TestFilterRolloverTasks(t *testing.T) {
	tests := []struct {
		name     string
		input    []Task
		expected []Task
	}{
		{
			name:     "empty slice",
			input:    []Task{},
			expected: []Task{},
		},
		{
			name: "all pending",
			input: []Task{
				{TaskTitle: "task1", Status: enums.Pending},
				{TaskTitle: "task2", Status: enums.Pending},
			},
			expected: []Task{
				{TaskTitle: "task1", Status: enums.Pending},
				{TaskTitle: "task2", Status: enums.Pending},
			},
		},
		{
			name: "all started",
			input: []Task{
				{TaskTitle: "task1", Status: enums.Started},
				{TaskTitle: "task2", Status: enums.Started},
			},
			expected: []Task{
				{TaskTitle: "task1", Status: enums.Started},
				{TaskTitle: "task2", Status: enums.Started},
			},
		},
		{
			name: "all complete",
			input: []Task{
				{TaskTitle: "task1", Status: enums.Complete},
				{TaskTitle: "task2", Status: enums.Complete},
			},
			expected: []Task{},
		},
		{
			name: "all abandoned",
			input: []Task{
				{TaskTitle: "task1", Status: enums.Abandoned},
				{TaskTitle: "task2", Status: enums.Abandoned},
			},
			expected: []Task{},
		},
		{
			name: "mixed statuses",
			input: []Task{
				{TaskTitle: "pending", Status: enums.Pending},
				{TaskTitle: "started", Status: enums.Started},
				{TaskTitle: "complete", Status: enums.Complete},
				{TaskTitle: "abandoned", Status: enums.Abandoned},
			},
			expected: []Task{
				{TaskTitle: "pending", Status: enums.Pending},
				{TaskTitle: "started", Status: enums.Started},
			},
		},
		{
			name: "preserves fields",
			input: []Task{
				{TaskTitle: "my task", TaskDesc: "my description", Status: enums.Pending},
			},
			expected: []Task{
				{TaskTitle: "my task", TaskDesc: "my description", Status: enums.Pending},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filterRolloverTasks(tt.input)

			if len(got) != len(tt.expected) {
				t.Fatalf("filterRolloverTasks() length = %d, want %d", len(got), len(tt.expected))
			}

			for i := range got {
				if got[i].TaskTitle != tt.expected[i].TaskTitle {
					t.Errorf("filterRolloverTasks() [%d].TaskTitle = %q, want %q", i, got[i].TaskTitle, tt.expected[i].TaskTitle)
				}

				if got[i].TaskDesc != tt.expected[i].TaskDesc {
					t.Errorf("filterRolloverTasks() [%d].TaskDesc = %q, want %q", i, got[i].TaskDesc, tt.expected[i].TaskDesc)
				}

				if got[i].Status != tt.expected[i].Status {
					t.Errorf("filterRolloverTasks() [%d].Status = %d, want %d", i, got[i].Status, tt.expected[i].Status)
				}
			}
		})
	}
}

func TestFindMostRecentFile(t *testing.T) {
	tests := []struct {
		name     string
		files    []string
		today    string
		expected string
	}{
		{
			name:     "empty directory",
			files:    []string{},
			today:    "2026-03-24",
			expected: "",
		},
		{
			name:     "returns yesterday",
			files:    []string{"2026-03-23"},
			today:    "2026-03-24",
			expected: "2026-03-23",
		},
		{
			name:     "skips today",
			files:    []string{"2026-03-24"},
			today:    "2026-03-24",
			expected: "",
		},
		{
			name:     "multi-day gap",
			files:    []string{"2026-03-18", "2026-03-20"},
			today:    "2026-03-24",
			expected: "2026-03-20",
		},
		{
			name:     "picks most recent before today",
			files:    []string{"2026-03-01", "2026-03-15", "2026-03-23"},
			today:    "2026-03-24",
			expected: "2026-03-23",
		},
		{
			name:     "ignores future files",
			files:    []string{"2026-03-20", "2026-03-25"},
			today:    "2026-03-24",
			expected: "2026-03-20",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			for _, f := range tt.files {
				os.WriteFile(filepath.Join(dir, f), []byte{}, 0o644)
			}

			got := findMostRecentFile(dir, tt.today)

			if tt.expected == "" {
				if got != "" {
					t.Errorf("findMostRecentFile() = %q, want empty string", got)
				}
			} else {
				want := filepath.Join(dir, tt.expected)
				if got != want {
					t.Errorf("findMostRecentFile() = %q, want %q", got, want)
				}
			}
		})
	}
}
