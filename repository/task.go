package repository

import (
	"a21hc3NpZ25tZW50/db/filebased"
	"a21hc3NpZ25tZW50/model"
	"fmt"
)

type TaskRepository interface {
	Store(task *model.Task) error
	Update(taskID int, task *model.Task) error
	Delete(id int) error
	GetByID(id int) (*model.Task, error)
	GetList() ([]model.Task, error)
	GetTaskCategory(id int) ([]model.TaskCategory, error)
}

type taskRepository struct {
	filebased *filebased.Data
}

func NewTaskRepo(filebasedDb *filebased.Data) *taskRepository {
	return &taskRepository{
		filebased: filebasedDb,
	}
}

// Store adds a new task to the database.
func (t *taskRepository) Store(task *model.Task) error {
	return t.filebased.StoreTask(*task)
}

// Update modifies an existing task in the database based on its ID.
func (t *taskRepository) Update(taskID int, task *model.Task) error {
	task.ID = taskID // ensure the task ID matches the one being updated
	return t.filebased.UpdateTask(taskID, *task)
}

// Delete removes a task from the database based on its ID.
func (t *taskRepository) Delete(id int) error {
	err := t.filebased.DeleteTask(id)
	if err != nil && err.Error() == "record not found" {
		// Return exactly "record not found" to match the test expectation.
		return fmt.Errorf("record not found")
	}
	return err
}

// GetByID retrieves a task by its ID from the database.
func (t *taskRepository) GetByID(id int) (*model.Task, error) {
	task, err := t.filebased.GetTaskByID(id)
	if err != nil && err.Error() == "record not found" {
		return nil, fmt.Errorf("record not found")
	}
	return task, err
}

// GetList retrieves all tasks from the database.
func (t *taskRepository) GetList() ([]model.Task, error) {
	tasks, err := t.filebased.GetTasks()
	if err != nil {
		return nil, fmt.Errorf("error fetching task list: %v", err)
	}
	return tasks, nil
}

// GetTaskCategory retrieves tasks with their category information based on category ID.
func (t *taskRepository) GetTaskCategory(id int) ([]model.TaskCategory, error) {
	taskCategories, err := t.filebased.GetTaskListByCategory(id)
	if err != nil {
		return nil, fmt.Errorf("error fetching tasks for category ID %d: %v", id, err)
	}
	return taskCategories, nil
}
