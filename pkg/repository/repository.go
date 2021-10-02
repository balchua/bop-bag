package repository

import (
	"database/sql"

	"github.com/balchua/bopbag/pkg/applog"
	"github.com/balchua/bopbag/pkg/domain"
	"go.uber.org/zap"
)

const (
	taskSchema = "CREATE TABLE IF NOT EXISTS TASKS (ID INTEGER PRIMARY KEY AUTOINCREMENT, TITLE VARCHAR(50), DETAILS VARCHAR(1000), CREATED_DATE VARCHAR(50), UNIQUE(ID))"
	insert     = "INSERT INTO TASKS (TITLE, DETAILS, CREATED_DATE) VALUES(?,?,?)"
	findById   = "SELECT ID, TITLE, DETAILS, CREATED_DATE FROM TASKS WHERE ID = ?"
	findAll    = "SELECT ID, TITLE, DETAILS, CREATED_DATE FROM TASKS"
)

type TaskRepository struct {
	db  Db
	log *applog.Logger
}

func NewTaskRepository(applog *applog.Logger, db Db) (*TaskRepository, error) {
	taskRepo := &TaskRepository{
		db:  db,
		log: applog,
	}
	err := taskRepo.migrate()

	if err != nil {
		return nil, err
	}
	return taskRepo, nil
}

func (t *TaskRepository) migrate() error {
	var err error
	if _, err = t.db.Exec(taskSchema); err != nil {
		t.log.Log.Fatal("unable to create schema", zap.Error(err))
	}
	return err
}

func (t *TaskRepository) Add(task *domain.Task) (*domain.Task, error) {
	var err error
	var result sql.Result
	if result, err = t.db.Exec(insert, task.Title, task.Details, task.CreatedDate); err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	t.log.Log.Info("record inserted", zap.Int64("records", rowsAffected))
	id, _ := result.LastInsertId()
	returnTask := &domain.Task{
		Id:          id,
		Title:       task.Title,
		Details:     task.Details,
		CreatedDate: task.CreatedDate,
	}
	return returnTask, err
}

func (t *TaskRepository) FindById(id int64) (*domain.Task, error) {
	var task domain.Task
	lg, _ := zap.NewProduction()
	lg.Info("Id to find", zap.Int64("id", id))

	row := t.db.QueryRow(findById, id)
	if err := row.Scan(&task.Id, &task.Title, &task.Details, &task.CreatedDate); err != nil {
		return nil, err
	}
	return &task, nil

}

func (t *TaskRepository) FindAll() (*[]domain.Task, error) {
	var tasks []domain.Task
	lg, _ := zap.NewProduction()

	var (
		taskId          sql.NullInt64
		taskTitle       sql.NullString
		taskDetails     sql.NullString
		taskCreatedDate sql.NullString
	)
	rows, err := t.db.Query(findAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var task domain.Task
		if err := rows.Scan(&taskId, &taskTitle, &taskDetails, &taskCreatedDate); err != nil {
			lg.Error("no record found", zap.Error(err))
			return nil, err
		}
		lg.Info("retrieved task", zap.Int64("id", taskId.Int64), zap.String("title", taskTitle.String),
			zap.String("details", taskDetails.String),
			zap.String("createdDate", taskCreatedDate.String))

		task = domain.Task{
			Id:          taskId.Int64,
			Title:       taskTitle.String,
			Details:     taskDetails.String,
			CreatedDate: taskCreatedDate.String,
		}

		tasks = append(tasks, task)
	}

	return &tasks, nil
}
