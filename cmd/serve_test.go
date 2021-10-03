package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/balchua/bopbag/pkg/applog"
	"github.com/stretchr/testify/assert"
)

func isOpened(host string, port int) bool {

	timeout := 5 * time.Second
	target := fmt.Sprintf("%s:%d", host, port)

	conn, err := net.DialTimeout("tcp", target, timeout)
	if err != nil {
		return false
	}

	if conn != nil {
		conn.Close()
		return true
	}

	return false
}

func TestStartApplication(t *testing.T) {
	applogger = applog.NewLogger()
	dir, err := ioutil.TempDir("/tmp/", "tempdb")
	if err != nil {
		log.Fatal(err)
	}
	enableTls = false
	dbAddress = "127.0.0.1:50000"
	dbPath = dir
	port = 8000
	startDqLite()
	defer os.Remove(dir)
	defer dqliteInst.CloseDqlite()

	startWiring()
	go startAppServer()

	time.Sleep(5 * time.Second)
	assert.True(t, isOpened("0.0.0.0", 8000))
}
