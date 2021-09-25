/*
Copyright © 2021 balchua

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/balchua/uncapsizable/pkg/controller"
	"github.com/balchua/uncapsizable/pkg/repository"
	"github.com/canonical/go-dqlite/app"
	"github.com/canonical/go-dqlite/client"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/sys/unix"
)

// serveCmd represents the serve command
var (
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Starts the application",
		Long:  `Starts the application along with its database`,
		Run:   start,
	}
	dbPath   string
	join     []string
	port     int
	dbPort   int
	taskRepo *repository.TaskRepository
)

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.PersistentFlags().StringVar(&dbPath, "db", "./", "Path to dqlite database files")
	serveCmd.PersistentFlags().StringSliceVar(&join, "join", []string{}, "Location of the main node to join to")
	serveCmd.PersistentFlags().IntVar(&port, "port", 8000, "Application web server port")
	serveCmd.PersistentFlags().IntVar(&dbPort, "dbPort", 9000, "the database port")

}

func startAppServer() {
	lg, _ := zap.NewProduction()

	taskController := controller.NewQueryController(taskRepo)

	// Fiber instance
	app := fiber.New()

	// Routes
	app.Get("/api/v1/task/:id", taskController.FindById)
	app.Get("/api/v1/tasks", taskController.FindAll)
	app.Post("/api/v1/task", taskController.NewTask)

	appErr := app.Listen(":" + strconv.Itoa(port))
	if appErr != nil {
		lg.Fatal("unable to start the app server")
	}
}

func dqliteLog(l client.LogLevel, format string, a ...interface{}) {
	log.Printf(fmt.Sprintf("%s\n", format), a...)
}

func startDqLite() {
	lg, _ := zap.NewProduction()
	var dqlite *app.App
	var err error
	if join == nil {
		dqlite, err = app.New(dbPath, app.WithAddress("0.0.0.0:"+strconv.Itoa(dbPort)), app.WithLogFunc(dqliteLog))
	} else {
		dqlite, err = app.New(dbPath, app.WithAddress("0.0.0.0:"+strconv.Itoa(dbPort)), app.WithCluster(join), app.WithLogFunc(dqliteLog))
	}

	if err != nil {
		lg.Fatal("Error while initializing dqlite %v", zap.Error(err))
	}
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	if err := dqlite.Ready(ctx); err != nil {
		lg.Fatal("Error while initializing dqlite %v", zap.Error(err))
	}
	db, _ := dqlite.Open(context.Background(), "uncapsizable")
	db.SetMaxOpenConns(5)
	db.SetConnMaxIdleTime(10 * time.Second)
	db.SetMaxIdleConns(5)
	lg.Info("database started")

	taskRepo, _ = repository.NewTaskRepository(db)
}

func start(cmd *cobra.Command, args []string) {

	startDqLite()
	startAppServer()

	ch := make(chan os.Signal)
	signal.Notify(ch, unix.SIGPWR)
	signal.Notify(ch, unix.SIGINT)
	signal.Notify(ch, unix.SIGQUIT)
	signal.Notify(ch, unix.SIGTERM)
	<-ch
}
