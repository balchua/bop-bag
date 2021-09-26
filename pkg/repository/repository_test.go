package repository

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/balchua/uncapsizable/pkg/applog"
	"github.com/balchua/uncapsizable/pkg/domain"
	"github.com/balchua/uncapsizable/pkg/infrastructure"
	"github.com/stretchr/testify/assert"
)

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
	dqliteInst, err := infrastructure.NewDqlite(applog, dir, dbAddress, nil)

	task := &domain.Task{
		Title:       "test",
		Details:     "test",
		CreatedDate: "20210926",
	}
	repo, err := NewTaskRepository(applog, dqliteInst)

	repo.Add(task)
	assert.NotNil(repo)
	assert.Nil(err)

	dqliteInst.CloseDqlite()
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
	dqliteInst, err := infrastructure.NewDqlite(applog, dir, dbAddress, nil)

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
	dqliteInst.CloseDqlite()
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
	dqliteInst, err := infrastructure.NewDqlite(applog, dir, dbAddress, nil)

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
	dqliteInst.CloseDqlite()
}
