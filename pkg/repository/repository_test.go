package repository

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/balchua/bopbag/pkg/applog"
	"github.com/balchua/bopbag/pkg/domain"
	"github.com/balchua/bopbag/pkg/infrastructure"
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

	dir, err := ioutil.TempDir("/tmp/", "tempdb")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(dir)
	dbAddress := "127.0.0.1:50000"
	dqliteInst, err := infrastructure.NewDqlite(applog, dir, dbAddress, nil, false, "")
	defer dqliteInst.Shutdown(context.TODO())
	task := &domain.Task{
		Title:       "test",
		Details:     "test",
		CreatedDate: "20210926",
	}
	repo, err := NewTaskRepository(applog, dqliteInst)

	repo.Add(task)
	assert.NotNil(repo)
	assert.Nil(err)
}

func TestSuccessfulFindOne(t *testing.T) {
	assert := assert.New(t)
	applog := applog.NewLogger()

	dir, err := ioutil.TempDir("/tmp/", "tempdb")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(dir)
	dbAddress := "127.0.0.1:50000"
	dqliteInst, err := infrastructure.NewDqlite(applog, dir, dbAddress, nil, false, "")
	defer dqliteInst.Shutdown(context.TODO())

	task := &domain.Task{
		Title:       "test",
		Details:     "test",
		CreatedDate: "20210926",
	}
	repo, err := NewTaskRepository(applog, dqliteInst)

	repo.Add(task)

	insertedTask, err := repo.FindById(int64(1))
	assert.NotNil(insertedTask)
	assert.Equal(insertedTask.Id, int64(1))
	assert.Nil(err)

}

func TestSuccessfulFindNothing(t *testing.T) {
	assert := assert.New(t)
	applog := applog.NewLogger()

	dir, err := ioutil.TempDir("/tmp/", "tempdb")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(dir)
	dbAddress := "127.0.0.1:50000"
	dqliteInst, err := infrastructure.NewDqlite(applog, dir, dbAddress, nil, false, "")
	defer dqliteInst.Shutdown(context.TODO())

	task := &domain.Task{
		Title:       "test",
		Details:     "test",
		CreatedDate: "20210926",
	}
	repo, err := NewTaskRepository(applog, dqliteInst)

	repo.Add(task)

	_, queryErr := repo.FindById(int64(2))
	assert.NotNil(t, queryErr)
	mustFail := ErrorContains(queryErr, "sql: no rows in result set")

	assert.True(mustFail)

}

func TestSuccessfulFindAll(t *testing.T) {
	assert := assert.New(t)
	applog := applog.NewLogger()

	dir, err := ioutil.TempDir("/tmp/", "tempdb")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(dir)
	dbAddress := "127.0.0.1:50000"
	dqliteInst, err := infrastructure.NewDqlite(applog, dir, dbAddress, nil, false, "")
	defer dqliteInst.Shutdown(context.TODO())

	task := &domain.Task{
		Title:       "test",
		Details:     "test",
		CreatedDate: "20210926",
	}
	repo, err := NewTaskRepository(applog, dqliteInst)

	repo.Add(task)

	tasks, err := repo.FindAll()
	assert.NotNil(tasks)
	assert.Equal(len(*tasks), 1)
	assert.Nil(err)

}

func TestSuccessfulFindAllWithMoreThanOneResult(t *testing.T) {
	assert := assert.New(t)
	applog := applog.NewLogger()

	dir, err := ioutil.TempDir("/tmp/", "tempdb")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(dir)
	dbAddress := "127.0.0.1:50000"
	dqliteInst, err := infrastructure.NewDqlite(applog, dir, dbAddress, nil, false, "")
	defer dqliteInst.Shutdown(context.TODO())
	repo, err := NewTaskRepository(applog, dqliteInst)

	for i := 0; i < 10; i++ {
		task := &domain.Task{
			Title:       "test" + strconv.Itoa(i),
			Details:     "test",
			CreatedDate: "20210926",
		}
		repo.Add(task)
	}

	tasks, err := repo.FindAll()
	assert.NotNil(tasks)
	assert.Equal(len(*tasks), 10)
	assert.Nil(err)

}

func TestSuccessfulFindAllEmptyDB(t *testing.T) {
	assert := assert.New(t)
	applog := applog.NewLogger()

	dir, err := ioutil.TempDir("/tmp/", "tempdb")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(dir)
	dbAddress := "127.0.0.1:50000"
	dqliteInst, err := infrastructure.NewDqlite(applog, dir, dbAddress, nil, false, "")
	defer dqliteInst.Shutdown(context.TODO())
	repo, err := NewTaskRepository(applog, dqliteInst)

	tasks, err := repo.FindAll()
	assert.Equal(len(*tasks), 0)
	assert.Nil(err)
}
