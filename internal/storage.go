package internal

type ListFilter int

const (
	FilterAll       ListFilter = 0
	FilterCompleted ListFilter = 1
	FilterPending   ListFilter = 2
)

type TaskAdder interface {
	AddTask(description string) error
}

type TaskCompleter interface {
	CompleteTask(id uint) error
}

type TaskLister interface {
	ListTasks(filter ListFilter) ([]*Task, error)
}

type TaskDeleter interface {
	DeleteTask(id uint) error
}

type Storage interface {
	TaskAdder
	TaskCompleter
	TaskLister
	TaskDeleter
}
