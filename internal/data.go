package internal

import "time"

type TaskStatus int

const (
	StatusPending   TaskStatus = 0
	StatusCompleted TaskStatus = 1
)

type Task struct {
	ID          uint
	Description string
	Status      TaskStatus
	DateCreated time.Time
	DateUpdated time.Time
}
