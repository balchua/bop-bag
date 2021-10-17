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
	"os"
	"os/signal"
	"strconv"

	"github.com/balchua/bopbag/pkg/applog"
	"github.com/balchua/bopbag/pkg/controller"
	"github.com/balchua/bopbag/pkg/infrastructure"
	"github.com/balchua/bopbag/pkg/repository"
	"github.com/balchua/bopbag/pkg/usecase"
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
	dbPath         string
	certsPath      string
	join           []string
	port           int
	dbAddress      string
	dqliteInst     *infrastructure.Dqlite
	taskRepo       *repository.TaskRepositoryImpl
	clusterRepo    *repository.ClusterRepository
	taskService    *usecase.TaskService
	taskController *controller.TaskController

	clusterController *controller.ClusterController
	clusterService    *usecase.ClusterService
	applogger         *applog.Logger
	enableTls         bool
)

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.PersistentFlags().StringVar(&dbPath, "db", "./", "Path to dqlite database files")
	serveCmd.PersistentFlags().StringSliceVar(&join, "join", []string{}, "Location of the main node to join to")
	serveCmd.PersistentFlags().IntVar(&port, "port", 8000, "Application web server port")
	serveCmd.PersistentFlags().StringVar(&dbAddress, "dbAddress", "localhost:9000", "the database port ex. localhost:9000")
	serveCmd.PersistentFlags().BoolVar(&enableTls, "enableTls", true, "Enable secure mode")
	serveCmd.PersistentFlags().StringVar(&certsPath, "certs", "./", "Path to dqlite certificates")

}

func startWiring() {
	retries := 5000
	taskRepo, _ = repository.NewTaskRepository(applogger, dqliteInst.DB())
	taskService = usecase.NewTaskService(taskRepo, uint(retries), applogger)
	clusterRepo = repository.NewClusterRepository(dqliteInst)
	clusterService = usecase.NewClusterService(clusterRepo, applogger)
	taskController = controller.NewTaskController(taskService)
	clusterController = controller.NewClusterController(clusterService)

}

func startAppServer() {
	// Fiber instance
	app := fiber.New()

	// Routes
	app.Get("/api/v1/task/:id", taskController.FindById)
	app.Get("/api/v1/tasks", taskController.FindAll)
	app.Post("/api/v1/task", taskController.NewTask)
	app.Put("/api/v1/task/:id", taskController.UpdateTask)
	app.Delete("/api/v1/task/:id", taskController.DeleteTask)
	app.Get("/api/v1/clusterInfo", clusterController.ShowCluster)
	app.Delete("/api/v1/node/:nodeId", clusterController.RemoveNode)

	appErr := app.Listen(":" + strconv.Itoa(port))
	if appErr != nil {
		applogger.Log.Fatal("unable to start the app server")
	}
}

func startDqLite() {
	var err error
	dqliteInst, err = infrastructure.NewDqlite(applogger, dbPath, dbAddress, join, enableTls, certsPath)

	if err != nil {
		applogger.Log.Fatal("unable to instantiate dqlite", zap.Error(err))
	}

}

func shutdownDqlite() {
	dqliteInst.Shutdown(context.Background())

}

func start(cmd *cobra.Command, args []string) {

	applogger = applog.NewLogger()
	startDqLite()
	startWiring()
	startAppServer()

	ch := make(chan os.Signal)
	signal.Notify(ch, unix.SIGPWR)
	signal.Notify(ch, unix.SIGINT)
	signal.Notify(ch, unix.SIGQUIT)
	signal.Notify(ch, unix.SIGTERM)
	<-ch

	shutdownDqlite()
}
