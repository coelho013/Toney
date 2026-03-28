package filetree

import "testing"

func TestSortPriority(t *testing.T) {
	tests := []struct {
		name     string
		node     Node
		expected int
	}{
		{"dot directory", Node{Name: ".git", IsDirectory: true}, 0},
		{"dot directory with longer name", Node{Name: ".obsidian", IsDirectory: true}, 0},
		{"regular directory", Node{Name: "journal", IsDirectory: true}, 1},
		{"dot file", Node{Name: ".gitignore", IsDirectory: false}, 2},
		{"regular file", Node{Name: "todo.md", IsDirectory: false}, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sortPriority(&tt.node)
			if got != tt.expected {
				t.Errorf("sortPriority(%q, IsDirectory=%v) = %d, want %d",
					tt.node.Name, tt.node.IsDirectory, got, tt.expected)
			}
		})
	}
}

func TestStartsWithDate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"YYYYMMDD", "20260328", true},
		{"YYYY-MM-DD", "2026-03-28", true},
		{"YYYYMMDD with suffix", "20260328-standup.md", true},
		{"YYYY-MM-DD with suffix", "2026-03-28-standup.md", true},
		{"regular file", "todo.md", false},
		{"short number", "2026", false},
		{"dot file", ".gitignore", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := startsWithDate(tt.input)
			if got != tt.expected {
				t.Errorf("startsWithDate(%q) = %v, want %v",
					tt.input, got, tt.expected)
			}
		})
	}
}
