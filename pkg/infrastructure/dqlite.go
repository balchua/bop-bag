package infrastructure

import (
	"context"
	"database/sql"
	"time"

	"github.com/balchua/uncapsizable/pkg/applog"
	"github.com/canonical/go-dqlite/app"
	"github.com/canonical/go-dqlite/client"
	"go.uber.org/zap"
)

const (
	MAX_CONNECTION           = 5
	MAX_IDLE_CONNECTION_TIME = 120 * time.Second
	DB_NAME                  = "uncapsizable"
)

type Dqlite struct {
	dqlite  *app.App
	address string
	log     *applog.Logger
	db      *sql.DB
}

func NewDqlite(log *applog.Logger, dbPath string, dbAddress string, join []string) (*Dqlite, error) {

	var dqlite *app.App
	var err error

	dqliteInstance := &Dqlite{}

	dqliteInstance.address = dbAddress
	dqliteInstance.log = log

	if join == nil {
		dqlite, err = app.New(dbPath, app.WithAddress(dqliteInstance.address), app.WithLogFunc(dqliteInstance.dqliteLog))
	} else {
		dqlite, err = app.New(dbPath, app.WithAddress(dqliteInstance.address), app.WithCluster(join), app.WithLogFunc(dqliteInstance.dqliteLog))
	}

	if err != nil {
		log.Log.Fatal("Error while initializing dqlite %v", zap.Error(err))
		return nil, err
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(300*time.Second))

	defer cancel()

	if err := dqlite.Ready(ctx); err != nil {
		log.Log.Sugar().Fatalf("Error while initializing dqlite %v", zap.Error(err))
		return nil, err
	}
	dqliteInstance.dqlite = dqlite
	if err := dqliteInstance.open(); err != nil {
		return nil, err
	}

	log.Log.Sugar().Infof("database %s started", DB_NAME)
	return dqliteInstance, nil
}

func (d *Dqlite) CloseDqlite() error {
	return d.dqlite.Close()
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

func (d *Dqlite) dqliteLog(l client.LogLevel, format string, a ...interface{}) {
	d.log.Log.Sugar().Info("[dqlite]", a)
}

func (d *Dqlite) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return d.db.Query(query, args...)
}

func (d *Dqlite) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return d.db.ExecContext(ctx, query, args...)
}
func (d *Dqlite) Exec(query string, args ...interface{}) (sql.Result, error) {
	return d.db.Exec(query, args...)
}

func (d *Dqlite) QueryRow(query string, args ...interface{}) *sql.Row {
	return d.db.QueryRow(query, args...)
}
