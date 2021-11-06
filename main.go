package main

import (
	"clean-arch/api"
	"clean-arch/csv"
	"clean-arch/infra"
	"clean-arch/todo"
	"clean-arch/worker"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	db := infra.NewMemoryDB()
	r := csv.NewCsvReader()
	uc := todo.NewUseCase(db, r)

	startWorker(uc)
	startServer(uc)
}

func startWorker(uc *todo.UseCase) {
	w := worker.NewWorker(uc)
	if err := w.ImportFromFile("todos.csv"); err != nil {
		log.Fatalln(err)
	}
}

func startServer(uc *todo.UseCase) {
	r := gin.Default()
	r.GET("/todos/pending", api.ListPendingTodosHandler(uc))
	r.PUT("/todo/:id/done", api.MarkTodoAsDoneHandler(uc))

	if err := r.Run(":8080"); err != nil {
		log.Fatalln(err)
	}
}
