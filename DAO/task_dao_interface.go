package DAO

import "J/model"

type TaskDAO interface {
	Create(task *model.Task) error
	GetList(filter model.TaskFilter) ([]*model.Task, error)
	Update(ID int, title string, done bool) error
	Delete(ID int) error
	Count(done bool) (int, error)
	Close() error
}
