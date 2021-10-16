package domain

type Task struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Details     string `json:"details"`
	CreatedDate string `json:"createdDate"`
}

type TaskRepository interface {
	Add(task *Task) (*Task, error)
	FindById(id int64) (*Task, error)
	FindAll() (*[]Task, error)
}
