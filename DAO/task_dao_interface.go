package DAO

import (
	"J/model"
	"time"
)

type TaskDAO interface {
	Create(task *model.Task) error
	GetList(filter model.TaskFilter) ([]*model.Task, error)
	Update(ID int, title string, done bool, ddl time.Time) error
	Delete(ID int) error
	Count(done bool) (int, error)
	Close() error
}
