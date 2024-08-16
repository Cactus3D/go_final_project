package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	addtask "cactus3d/go_final_project/internal/http-server/handlers/add-task"
	deletetask "cactus3d/go_final_project/internal/http-server/handlers/delete-task"
	donetask "cactus3d/go_final_project/internal/http-server/handlers/done-task"
	gettask "cactus3d/go_final_project/internal/http-server/handlers/get-task"
	gettasks "cactus3d/go_final_project/internal/http-server/handlers/get-tasks"
	nextdate "cactus3d/go_final_project/internal/http-server/handlers/next-date"
	updatetask "cactus3d/go_final_project/internal/http-server/handlers/update-task"
	tasks "cactus3d/go_final_project/internal/service/tasks"
	sqlite "cactus3d/go_final_project/internal/storage/sqlite"

	"github.com/go-chi/chi/v5"
)

const (
	defaultWebDir = "./web/"
	defaultDBFile = "./scheduler.db"
	defaultPort   = 7540
)

func main() {

	dbFile := os.Getenv("TODO_DBFILE")

	if dbFile == "" {
		dbFile = defaultDBFile
	}
	store, err := sqlite.New(dbFile)
	if err != nil {
		log.Fatalf("error starting db: %v", err)
		return
	}
	defer func() {
		err = store.Close()
		if err != nil {
			log.Fatalf("DB failed to shutdown: %v\n", err)
		}
	}()

	portStr := os.Getenv("TODO_PORT")
	var port int

	if portStr == "" {
		port = defaultPort
	} else {
		port, err = strconv.Atoi(portStr)
		if err != nil {
			log.Fatalf("Invalid port number: %v", err)
			return
		}
	}

	webDir := os.Getenv("TODO_WEB_DIR")

	if webDir == "" {
		webDir = defaultWebDir
	}

	taskService := tasks.New(store)

	r := chi.NewRouter()

	r.Handle("/*", http.FileServer(http.Dir(webDir)))

	r.MethodFunc(http.MethodGet, "/api/nextdate", nextdate.New())

	r.MethodFunc(http.MethodPost, "/api/task", addtask.New(taskService))
	r.MethodFunc(http.MethodGet, "/api/task", gettask.New(taskService))
	r.MethodFunc(http.MethodPut, "/api/task", updatetask.New(taskService))
	r.MethodFunc(http.MethodDelete, "/api/task", deletetask.New(taskService))

	r.MethodFunc(http.MethodGet, "/api/tasks", gettasks.New(taskService))

	r.MethodFunc(http.MethodPost, "/api/task/done", donetask.New(taskService))

	address := fmt.Sprintf(":%d", port)
	if err := http.ListenAndServe(address, r); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed with err: %v\n", err)
		return
	}
}
