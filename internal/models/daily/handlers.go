package daily

import (
	"slices"

	"github.com/SourcewareLab/Toney/v2/internal/enums"
	"github.com/SourcewareLab/Toney/v2/internal/messages"
)

func (m Daily) CreateTask(msg messages.TaskPopupMessage, isUnique bool) error {
	task := Task{
		TaskTitle: msg.Title,
		TaskDesc:  msg.Desc,
		Status:    msg.Status,
	}

	if isUnique {
		task.TaskType = enums.UniqueTask
	} else {
		task.TaskType = enums.RecurringTask
	}

	m.Tasks = append(m.Tasks, task)

	return WriteItems(m.Tasks)
}

func (m Daily) DeleteTask(msg messages.TaskPopupMessage) error {
	if !msg.IsDeleted {
		return nil
	}

	item := m.List.SelectedItem()

	task, ok := item.(Task)
	if !ok { // Making sure that item is of type Task
		return nil
	}

	m.Tasks = slices.Delete(m.Tasks, task.ID-1, task.ID)

	return WriteItems(m.Tasks)
}

func (m Daily) StatusChangeTask(msg messages.TaskPopupMessage) error {
	item := m.List.SelectedItem()

	task, ok := item.(Task)
	if !ok { // Making sure that item is of type Task
		return nil
	}

	task.Status = msg.Status

	m.Tasks[task.ID-1] = task

	return WriteItems(m.Tasks)
}

func (m Daily) EditTask(msg messages.TaskPopupMessage) error {
	task := Task{
		TaskTitle: msg.Title,
		TaskDesc:  msg.Desc,
		Status:    msg.Status,
	}

	item := m.List.SelectedItem()

	oldTask, ok := item.(Task)
	if !ok { // Making sure that item is of type Task
		return nil
	}

	task.Status = oldTask.Status

	m.Tasks[oldTask.ID-1] = task

	return WriteItems(m.Tasks)
}
