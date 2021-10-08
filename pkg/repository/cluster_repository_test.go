package repository

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/balchua/bopbag/pkg/applog"
	"github.com/balchua/bopbag/pkg/infrastructure"
	"github.com/stretchr/testify/assert"
)

func TestMustRetrieveClusterInfo(t *testing.T) {
	assert := assert.New(t)
	applog := applog.NewLogger()

	dir, err := ioutil.TempDir("/tmp/", "tempdb")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(dir)
	dbAddress := "127.0.0.1:50000"
	dqliteInst, err := infrastructure.NewDqlite(applog, dir, dbAddress, nil, false, "")
	defer dqliteInst.CloseDqlite()

	repo := NewClusterRepository(dqliteInst)

	data, err := repo.ClusterInfo()
	assert.NotNil(data)

	var indented bytes.Buffer
	json.Indent(&indented, data, "", "\t")
	result := string(indented.Bytes())
	assert.Contains(result, "127.0.0.1:50000")
}
