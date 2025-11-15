package model

import "time"

type Task struct {
	ID       int       `json:"id"`
	Title    string    `json:"title"`
	Done     bool      `json:"done"`
	CreateAt time.Time `json:"create_at"`
	UpdateAt time.Time `json:"update_at"`
	DeadLine time.Time `json:"dead_line"`
}

type TaskFilter struct {
	Done            bool
	Limit           int
	OrderByDeadline bool
}
