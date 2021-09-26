package domain

type Task struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Details     string `json:"details"`
	CreatedDate string `json:"createdDate"`
}
