package daily

import (
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/SourcewareLab/Toney/v2/internal/config"
	"github.com/SourcewareLab/Toney/v2/internal/enums"
	"github.com/charmbracelet/bubbles/list"
	"github.com/jszwec/csvutil"
)

type Task struct {
	TaskTitle string           `csv:"title"`
	TaskDesc  string           `csv:"desc"`
	Status    enums.TaskStatus `csv:"status"`
	// We will not store these data in the file
	ID       int            `csv:"-"` // Point to index in the respective type array
	TaskType enums.TaskType `csv:"-"`
}

func (m Task) Title() string       { return m.TaskTitle }
func (m Task) Description() string { return m.TaskDesc }
func (m Task) FilterValue() string { return m.TaskTitle }

func TaskToItems(tasks []Task) []list.Item {
	list := make([]list.Item, 0)
	for i, v := range tasks {
		v.ID = i + 1
		v.TaskType = enums.RecurringTask
		list = append(list, v)
	}
	return list
}

func GetItems() ([]Task, error) {
	path := GetPath()

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		f, err2 := os.Create(path)
		if err2 != nil {
			return nil, err2
		}

		prevPath := getMostRecentDailyPath()
		if prevPath != "" {
			content, err2 := os.ReadFile(prevPath)
			if err2 != nil {
				return nil, err2
			}

			var allTasks []Task
			if err := csvutil.Unmarshal(content, &allTasks); err != nil {
				return nil, err
			}

			tasks := filterRolloverTasks(allTasks)

			data, err2 := csvutil.Marshal(tasks)
			if err2 != nil {
				return nil, err2
			}

			if _, err := f.Write(data); err != nil {
				return nil, err
			}
		}
	} else if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	tasks := make([]Task, 0)
	err = csvutil.Unmarshal(content, &tasks)
	if err != nil {
		return nil, err
	}

	slices.SortStableFunc(tasks, func(a, b Task) int {
		return getStatusOrder(a.Status) - getStatusOrder(b.Status)
	})

	return tasks, nil
}

func WriteItems(tasks []Task) error {
	path := GetPath()

	data, err := csvutil.Marshal(tasks)
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return err
	}

	return nil
}

func GetPath() string {
	home, _ := os.UserHomeDir()
	date := time.Now().Format("2006-01-02")

	return filepath.Join(home, config.AppConfig.General.NotesDir, ".daily", date)
}

func getMostRecentDailyPath() string {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, config.AppConfig.General.NotesDir, ".daily")
	today := time.Now().Format("2006-01-02")

	return findMostRecentFile(dir, today)
}

func findMostRecentFile(dir string, today string) string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}

	for i := len(entries) - 1; i >= 0; i-- {
		name := entries[i].Name()
		if !entries[i].IsDir() && name < today {
			return filepath.Join(dir, name)
		}
	}

	return ""
}

func filterRolloverTasks(tasks []Task) []Task {
	result := make([]Task, 0)
	for _, task := range tasks {
		if task.Status == enums.Complete || task.Status == enums.Abandoned {
			continue
		}

		result = append(result, task)
	}
	return result
}

func getStatusOrder(s enums.TaskStatus) int {
	switch s {
	case enums.Started:
		return 0
	case enums.Pending:
		return 1
	case enums.Abandoned:
		return 2
	case enums.Complete:
		return 3
	default:
		return 4
	}
}
