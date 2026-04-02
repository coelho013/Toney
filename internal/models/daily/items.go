package daily

import (
	"fmt"
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

func GetItems() []Task {
	path := GetPath()

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		f, err2 := os.Create(path)
		if err2 != nil {
			fmt.Println(err2.Error())
		}

		prevPath := getMostRecentDailyPath()
		if prevPath != "" {
			content, err2 := os.ReadFile(prevPath)
			if err2 != nil {
				fmt.Println(err2.Error())
			}

			var allTasks []Task
			csvutil.Unmarshal(content, &allTasks)

			tasks := filterRolloverTasks(allTasks)

			data, err2 := csvutil.Marshal(tasks)
			if err2 != nil {
				fmt.Println(err2.Error())
			}

			f.Write(data)
		}
	} else if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	content, _ := os.ReadFile(path)

	tasks := make([]Task, 0)
	err = csvutil.Unmarshal(content, &tasks)
	if err != nil {
		fmt.Println(err.Error())
	}

	slices.SortStableFunc(tasks, func(a, b Task) int {
		return getStatusOrder(a.Status) - getStatusOrder(b.Status)
	})

	return tasks
}

func WriteItems(tasks []Task) {
	path := GetPath()

	data, _ := csvutil.Marshal(tasks)

	os.WriteFile(path, data, 0o644)
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
