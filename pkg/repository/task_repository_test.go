package repository

import (
	"fmt"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/balchua/bopbag/pkg/applog"
	"github.com/balchua/bopbag/pkg/domain"
	"github.com/stretchr/testify/assert"
)

func ErrorContains(out error, want string) bool {
	if out == nil {
		return want == ""
	}
	if want == "" {
		return false
	}
	return strings.Contains(out.Error(), want)
}

// TestSuccessfulInsert is an integration test, it runs a real Dqlite instance and delete it afterwards
func TestSuccessfulInsert(t *testing.T) {
	assert := assert.New(t)
	applog := applog.NewLogger()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	task := &domain.Task{
		Title:       "test",
		Details:     "test",
		CreatedDate: "20210926",
	}
	mock.ExpectExec("INSERT INTO TASKS").WithArgs("test", "test", "20210926").WillReturnResult(sqlmock.NewResult(1, 1))
	repo, err := NewTaskRepository(applog, db)
	repo.Add(task)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	assert.NotNil(repo)
	assert.Nil(err)
}

func TestFailInsert(t *testing.T) {
	assert := assert.New(t)
	applog := applog.NewLogger()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	task := &domain.Task{
		Title:       "test",
		Details:     "test",
		CreatedDate: "20210926",
	}
	mock.ExpectExec("INSERT INTO TASKS").WithArgs("test", "test", "20210926").
		WillReturnError(fmt.Errorf("database error"))

	repo, err := NewTaskRepository(applog, db)
	_, addErr := repo.Add(task)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	assert.NotNil(addErr)
}

func TestSuccessfulFindOne(t *testing.T) {
	assert := assert.New(t)
	applog := applog.NewLogger()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	columns := []string{"id", "title", "details", "created_date"}
	row := sqlmock.NewRows(columns).AddRow(int64(1), "test", "test", "20210926")
	mock.ExpectQuery("SELECT ID, TITLE, DETAILS, CREATED_DATE FROM TASKS WHERE ID = ?").
		WithArgs(int64(1)).
		WillReturnRows(row)

	repo, err := NewTaskRepository(applog, db)

	insertedTask, err := repo.FindById(int64(1))
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	assert.NotNil(insertedTask)
	assert.Equal(insertedTask.Id, int64(1))
	assert.Equal(insertedTask.Title, "test")
	assert.Equal(insertedTask.Details, "test")
	assert.Equal(insertedTask.CreatedDate, "20210926")
	assert.Nil(err)

}

func TestFailFindOne(t *testing.T) {
	assert := assert.New(t)
	applog := applog.NewLogger()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	mock.ExpectQuery("SELECT ID, TITLE, DETAILS, CREATED_DATE FROM TASKS WHERE ID = ?").
		WillReturnError(fmt.Errorf("database error"))

	repo, err := NewTaskRepository(applog, db)

	_, findErr := repo.FindById(int64(1))
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	assert.NotNil(findErr)

}

func TestSuccessfulFindNothing(t *testing.T) {
	assert := assert.New(t)
	applog := applog.NewLogger()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	columns := []string{"id", "title", "details", "created_date"}
	emptyResult := sqlmock.NewRows(columns)
	mock.ExpectQuery("SELECT ID, TITLE, DETAILS, CREATED_DATE FROM TASKS WHERE ID = ?").
		WithArgs(int64(2)).
		WillReturnRows(emptyResult)

	repo, err := NewTaskRepository(applog, db)

	_, queryErr := repo.FindById(int64(2))

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	assert.NotNil(t, queryErr)
	mustFail := ErrorContains(queryErr, "sql: no rows in result set")

	assert.True(mustFail)

}

func TestSuccessfulFindAll(t *testing.T) {
	assert := assert.New(t)
	applog := applog.NewLogger()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	columns := []string{"id", "title", "details", "created_date"}
	row := sqlmock.NewRows(columns).AddRow(int64(1), "test", "test", "20210926").AddRow(int64(2), "test2", "test2", "20211019")
	mock.ExpectQuery("SELECT ID, TITLE, DETAILS, CREATED_DATE FROM TASKS").
		WillReturnRows(row)

	repo, err := NewTaskRepository(applog, db)

	tasks, err := repo.FindAll()

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	assert.NotNil(tasks)
	assert.Equal(len(*tasks), 2)
	assert.Nil(err)

}

func TestSuccessfulFindAllEmptyDB(t *testing.T) {
	assert := assert.New(t)
	applog := applog.NewLogger()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	columns := []string{"id", "title", "details", "created_date"}
	emptyResult := sqlmock.NewRows(columns)
	mock.ExpectQuery("SELECT ID, TITLE, DETAILS, CREATED_DATE FROM TASKS").
		WillReturnRows(emptyResult)

	repo, err := NewTaskRepository(applog, db)

	tasks, err := repo.FindAll()
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.Equal(len(*tasks), 0)
	assert.Nil(err)
}

func TestFailedFindAll(t *testing.T) {
	assert := assert.New(t)
	applog := applog.NewLogger()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	mock.ExpectQuery("SELECT ID, TITLE, DETAILS, CREATED_DATE FROM TASKS").
		WillReturnError(fmt.Errorf("database error"))

	repo, err := NewTaskRepository(applog, db)

	_, findAllErr := repo.FindAll()
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(findAllErr)
}

func TestSuccessfulDeleteTask(t *testing.T) {
	id := int64(1)
	assert := assert.New(t)
	applog := applog.NewLogger()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectExec("DELETE FROM TASKS WHERE ID = ?").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo, err := NewTaskRepository(applog, db)

	deleteErr := repo.Delete(id)
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	assert.Nil(deleteErr)
}

func TestFailedDelete(t *testing.T) {
	id := int64(1)
	assert := assert.New(t)
	applog := applog.NewLogger()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectExec("DELETE FROM TASKS WHERE ID = ?").
		WithArgs(id).
		WillReturnError(fmt.Errorf("database error"))

	repo, err := NewTaskRepository(applog, db)

	deleteErr := repo.Delete(id)
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	assert.NotNil(deleteErr)
}

func TestSuccessfulUpdate(t *testing.T) {
	id := int64(1)
	assert := assert.New(t)
	applog := applog.NewLogger()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	task := &domain.Task{
		Id:          id,
		Title:       "update title",
		Details:     "new details",
		CreatedDate: "20210926",
	}
	mock.ExpectExec("UPDATE TASKS ").
		WithArgs("update title", "new details", id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo, err := NewTaskRepository(applog, db)
	updatedTask, updateErr := repo.Update(task)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	assert.NotNil(updatedTask)
	assert.Equal(updatedTask.Id, int64(1))
	assert.Equal(updatedTask.Title, "update title")
	assert.Equal(updatedTask.Details, "new details")
	assert.Nil(updateErr)
}

func TestFailedUpdate(t *testing.T) {
	id := int64(1)
	assert := assert.New(t)
	applog := applog.NewLogger()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	task := &domain.Task{
		Id:          id,
		Title:       "update title",
		Details:     "new details",
		CreatedDate: "20210926",
	}
	mock.ExpectExec("UPDATE TASKS ").
		WithArgs("update title", "new details", id).
		WillReturnError(fmt.Errorf("database error"))

	repo, err := NewTaskRepository(applog, db)
	_, updateErr := repo.Update(task)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	assert.NotNil(updateErr)
}
