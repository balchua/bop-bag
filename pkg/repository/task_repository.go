package repository

import (
	"database/sql"
	"time"

	"github.com/balchua/bopbag/pkg/applog"
	"github.com/balchua/bopbag/pkg/domain"
	"go.uber.org/zap"
)

const (
	insert   = "INSERT INTO TASKS (TITLE, DETAILS, CREATED_DATE) VALUES(?,?,?)"
	delete   = "DELETE FROM TASKS WHERE ID = ?"
	findById = "SELECT ID, TITLE, DETAILS, CREATED_DATE FROM TASKS WHERE ID = ?"
	findAll  = "SELECT ID, TITLE, DETAILS, CREATED_DATE FROM TASKS"
	update   = "UPDATE TASKS SET TITLE=?, DETAILS=? WHERE ID=?"
)

type TaskRepositoryImpl struct {
	db  *sql.DB
	log *applog.Logger
}

func NewTaskRepository(applog *applog.Logger, db *sql.DB) (*TaskRepositoryImpl, error) {
	taskRepo := &TaskRepositoryImpl{
		db:  db,
		log: applog,
	}
	return taskRepo, nil
}

func (t *TaskRepositoryImpl) Add(task *domain.Task) (*domain.Task, error) {
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

func (t *TaskRepositoryImpl) FindById(id int64) (*domain.Task, error) {
	var task domain.Task
	lg, _ := zap.NewProduction()
	lg.Info("Id to find", zap.Int64("id", id))

	row := t.db.QueryRow(findById, id)
	if err := row.Scan(&task.Id, &task.Title, &task.Details, &task.CreatedDate); err != nil {
		return nil, err
	}
	return &task, nil

}

func (t *TaskRepositoryImpl) FindAll() (*[]domain.Task, error) {
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

func (t *TaskRepositoryImpl) Delete(id int64) error {
	var err error
	var result sql.Result
	var rowsAffected int64

	lg, _ := zap.NewProduction()

	lg.Info("Id to find", zap.Int64("id", id))

	if result, err = t.db.Exec(delete, id); err != nil {
		return err
	}
	rowsAffected, err = result.RowsAffected()

	lg.Info("number of rows affected", zap.Int64("rows", rowsAffected))

	return nil

}

func (t *TaskRepositoryImpl) Update(task *domain.Task) (*domain.Task, error) {
	var err error
	var result sql.Result
	currentTime := time.Now()
	task.CreatedDate = currentTime.Format(time.RFC1123)
	if result, err = t.db.Exec(update, task.Title, task.Details, task.Id); err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	t.log.Log.Info("record updated", zap.Int64("records", rowsAffected))
	id := task.Id
	returnTask := &domain.Task{
		Id:          id,
		Title:       task.Title,
		Details:     task.Details,
		CreatedDate: task.CreatedDate,
	}
	return returnTask, err
}
