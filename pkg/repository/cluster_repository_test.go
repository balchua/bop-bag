package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/balchua/bopbag/pkg/applog"
	"github.com/balchua/bopbag/pkg/infrastructure"
	_ "github.com/balchua/bopbag/pkg/test_util"
	"github.com/stretchr/testify/assert"
)

func TestMustRetrieveClusterInfo(t *testing.T) {
	assert := assert.New(t)
	applog := applog.NewLogger()
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	certsPath := path + "/default-certs/"
	enableTls := true
	dir, err := ioutil.TempDir("/tmp/", "tempdb")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(dir)
	dbAddress := "127.0.0.1:50000"
	dqliteInst, err := infrastructure.NewDqlite(applog, dir, dbAddress, nil, enableTls, certsPath)
	defer dqliteInst.Shutdown(context.TODO())

	repo := NewClusterRepository(dqliteInst)

	data, err := repo.ClusterInfo()
	assert.NotNil(data)

	var indented bytes.Buffer
	json.Indent(&indented, data, "", "\t")
	result := string(indented.Bytes())
	assert.Contains(result, "127.0.0.1:50000")
}

func TestTryRemoveNode(t *testing.T) {
	assert := assert.New(t)
	applog := applog.NewLogger()
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	certsPath := path + "/default-certs/"
	enableTls := true
	dir, err := ioutil.TempDir("/tmp/", "tempdb")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(dir)
	dbAddress := "127.0.0.1:50000"
	dqliteInst, err := infrastructure.NewDqlite(applog, dir, dbAddress, nil, enableTls, certsPath)
	defer dqliteInst.Shutdown(context.TODO())

	repo := NewClusterRepository(dqliteInst)

	data, err := repo.RemoveNode(dbAddress)
	assert.Nil(err)
	assert.Equal("127.0.0.1:50000", data)
}

func TestRemoveNoneExistentNode(t *testing.T) {
	assert := assert.New(t)
	applog := applog.NewLogger()
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	certsPath := path + "/default-certs/"
	enableTls := true
	dir, err := ioutil.TempDir("/tmp/", "tempdb")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(dir)
	dbAddress := "127.0.0.1:50000"
	dqliteInst, err := infrastructure.NewDqlite(applog, dir, dbAddress, nil, enableTls, certsPath)
	defer dqliteInst.Shutdown(context.TODO())

	repo := NewClusterRepository(dqliteInst)

	_, removeErr := repo.RemoveNode("127.0.0.1:2000")
	assert.NotNil(removeErr)
}

func TestFailDbStartupNoDataDirectory(t *testing.T) {
	assert := assert.New(t)
	applog := applog.NewLogger()

	dbAddress := "127.0.0.1:50000"
	_, err := infrastructure.NewDqlite(applog, "/non-existent/", dbAddress, nil, false, "")

	assert.NotNil(t, err)
}

func TestGetLeader(t *testing.T) {
	assert := assert.New(t)
	applog := applog.NewLogger()
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	certsPath := path + "/default-certs/"
	enableTls := true
	dir, err := ioutil.TempDir("/tmp/", "tempdb")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(dir)
	dbAddress := "127.0.0.1:50000"
	dqliteInst, err := infrastructure.NewDqlite(applog, dir, dbAddress, nil, enableTls, certsPath)
	defer dqliteInst.Shutdown(context.TODO())

	repo := NewClusterRepository(dqliteInst)

	data, err := repo.FindLeader()
	assert.Nil(err)
	assert.Equal("127.0.0.1:50000", data)
}
