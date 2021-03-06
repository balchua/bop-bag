package infrastructure

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/balchua/bopbag/pkg/applog"
	"github.com/canonical/go-dqlite/app"
	"github.com/canonical/go-dqlite/client"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	MAX_CONNECTION           = 2
	MAX_IDLE_CONNECTION_TIME = 2 * time.Second
	DB_NAME                  = "bopbag"
	taskSchema               = "CREATE TABLE IF NOT EXISTS TASKS (ID INTEGER PRIMARY KEY AUTOINCREMENT, TITLE VARCHAR(50), DETAILS VARCHAR(1000), CREATED_DATE VARCHAR(50), UNIQUE(ID))"
)

type Dqlite struct {
	dqlite  *app.App
	address string
	log     *applog.Logger
	db      *sql.DB
}

func NewDqlite(log *applog.Logger, dbPath string, dbAddress string, join []string, enableTls bool, certsPath string) (*Dqlite, error) {

	var dqlite *app.App
	var err error

	dqliteInstance := &Dqlite{}

	dqliteInstance.address = dbAddress
	dqliteInstance.log = log

	options := []app.Option{
		app.WithAddress(dqliteInstance.address),
		app.WithLogFunc(dqliteInstance.dqliteLog),
	}

	if enableTls {
		// Load the TLS certificates.
		crt := filepath.Join(certsPath, "cluster.crt")
		key := filepath.Join(certsPath, "cluster.key")

		keypair, err := tls.LoadX509KeyPair(crt, key)
		if err != nil {
			return nil, errors.Wrap(err, "load keypair")
		}
		data, err := ioutil.ReadFile(crt)
		if err != nil {
			return nil, errors.Wrap(err, "read certificate")
		}
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(data) {
			return nil, fmt.Errorf("bad certificate")
		}
		options = append(options, app.WithTLS(app.SimpleTLSConfig(keypair, pool)))
	}

	if join != nil {
		options = append(options, app.WithCluster(join))
	}

	dqlite, err = app.New(dbPath, options...)

	if err != nil {
		log.Log.Sugar().Errorf("Error while initializing dqlite %v", zap.Error(err))
		return nil, err
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(120*time.Second))

	defer cancel()

	if err := dqlite.Ready(ctx); err != nil {
		log.Log.Sugar().Errorf("Error while initializing dqlite %v", zap.Error(err))
		return nil, err
	}
	dqliteInstance.dqlite = dqlite
	if err := dqliteInstance.open(); err != nil {
		return nil, err
	}

	if err := dqliteInstance.migrate(); err != nil {
		return nil, err
	}

	log.Log.Sugar().Infof("database %s started", DB_NAME)
	return dqliteInstance, nil
}

func (d *Dqlite) open() error {
	db, err := d.dqlite.Open(context.Background(), DB_NAME)
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(MAX_CONNECTION)
	db.SetConnMaxIdleTime(MAX_IDLE_CONNECTION_TIME)
	db.SetMaxIdleConns(MAX_CONNECTION)
	d.db = db

	return nil
}

func (d *Dqlite) migrate() error {
	var err error
	if _, err = d.db.Exec(taskSchema); err != nil {
		d.log.Log.Fatal("unable to create schema", zap.Error(err))
	}
	return err
}

func (d *Dqlite) dqliteLog(l client.LogLevel, format string, a ...interface{}) {
	d.log.Log.Sugar().Infof("[dqlite] %s - %v", format, a)
}

func (d *Dqlite) DB() *sql.DB {
	return d.db
}

func (d *Dqlite) GetClusterInfo() ([]byte, error) {
	ctx := context.Background()
	cli, err := d.dqlite.Client(ctx)

	if err != nil {
		return nil, err
	}
	cluster, err := cli.Cluster(ctx)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(cluster)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (d *Dqlite) RemoveNode(address string) (string, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))

	defer cancel()
	cli, err := d.dqlite.Leader(ctx)
	if err != nil {
		return "", err
	}
	cluster, err := cli.Cluster(ctx)
	if err != nil {
		return "", err
	}

	for _, node := range cluster {
		if node.Address != address {
			continue
		}
		//ignore error returned by dqlite
		err := cli.Remove(ctx, node.ID)
		if err != nil {
			d.log.Log.Sugar().Errorf("Error while removing a node %v", zap.Error(err))
		}
		return address, nil
	}
	return "", fmt.Errorf("no node has address %q", address)
}

func (d *Dqlite) Leader() (string, error) {
	var leader *client.NodeInfo
	var err error
	var cli *client.Client
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))

	defer cancel()

	cli, err = d.dqlite.Leader(ctx)
	if err != nil {
		return "", err
	}

	leader, err = cli.Leader(ctx)
	if err != nil {
		return "", err
	}

	d.log.Log.Sugar().Infof("Current Leader %s", leader.Address)
	return leader.Address, nil

}

func (d *Dqlite) Shutdown(ctx context.Context) {
	if err := d.db.Close(); err != nil {
		d.log.Log.Sugar().Errorf("Unable to close the db %v", err)
	}
	if err := d.dqlite.Handover(ctx); err != nil {
		d.log.Log.Sugar().Errorf("Unable to handover leadership %v", err)
	}
	if err := d.dqlite.Close(); err != nil {
		d.log.Log.Sugar().Errorf("Unable to shutdown dqlite %v", err)
	}
}
