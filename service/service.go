package service

import (
	"J/DAO"
	"J/model"
	"fmt"
	"time"
)

type TodoService struct {
	taskDAO DAO.TaskDAO
}

func NewTodoService(taskDAO DAO.TaskDAO) *TodoService {
	return &TodoService{
		taskDAO: taskDAO,
	}
}
func (service *TodoService) AddTask(title string, ddl int) (*model.Task, error) {
	if title == "" {
		return nil, fmt.Errorf("task title can not be empty")
	}
	task := &model.Task{
		Title:    title,
		Done:     false,
		DeadLine: time.Now().Add(time.Duration(time.Minute * time.Duration(ddl))),
	}
	err := service.taskDAO.Create(task)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %v", err)
	}
	return task, nil
}

func (service *TodoService) ShowUndoTasks() ([]*model.Task, error) {
	undoFilter := model.TaskFilter{Done: false}
	undoTasks, err := service.taskDAO.GetList(undoFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to get undo tasks: %v", err)
	}
	return undoTasks, nil
}

func (service *TodoService) ShowDoneTasks() ([]*model.Task, error) {
	doneFilter := model.TaskFilter{Done: true}
	doneTasks, err := service.taskDAO.GetList(doneFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to get done tasks: %v", err)
	}
	return doneTasks, nil
}

func (service *TodoService) UpdateTask(ID int, title string, done bool, ddl time.Time) error {
	err := service.taskDAO.Update(ID, title, done, ddl)
	if err != nil {
		return fmt.Errorf("failed to update task: %v", err)
	}
	return nil
}

func (service *TodoService) DeleteTask(ID int) error {
	err := service.taskDAO.Delete(ID)
	if err != nil {
		return fmt.Errorf("failed to delete task: %v", err)
	}
	return nil
}

// new
func (service *TodoService) GetUrgentTasks(limit int) ([]*model.Task, error) {
	filter := model.TaskFilter{
		Done:            false,
		Limit:           limit,
		OrderByDeadline: true,
	}
	tasks, err := service.taskDAO.GetList(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get urgent tasks:%v", err)
	}
	return tasks, nil
}

func (service *TodoService) FinishedTask(ID int) error {
	done := true
	return service.UpdateTask(ID, "", done, time.Time{})
}

func (service *TodoService) ClearAllTasks() error {
	Tasks := make([]*model.Task, 0)
	UndoTasks, err := service.ShowUndoTasks()
	if err != nil {
		return fmt.Errorf("failed to clear all tasks: %v", err)
	}
	Tasks = append(Tasks, UndoTasks...)
	DoneTasks, err := service.ShowDoneTasks()
	if err != nil {
		return fmt.Errorf("failed to clear all tasks: %v", err)
	}
	Tasks = append(Tasks, DoneTasks...)
	for _, task := range Tasks {
		err := service.taskDAO.Delete(task.ID)
		if err != nil {
			fmt.Printf("ID: %d delete error: %v", task.ID, err)
		}
	}

	return nil
}

func (service *TodoService) GetRecentUndoTasks(limit int) ([]*model.Task, error) {
	undoTasks, err := service.ShowUndoTasks()
	if err != nil {
		return nil, err
	}
	if limit > 0 && limit < len(undoTasks) {
		undoTasks = undoTasks[:limit]
	}
	return undoTasks, nil
}

func (service *TodoService) Close() error {
	return service.taskDAO.Close()
}
