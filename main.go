package main

import (
	"clean-arch/api"
	"clean-arch/csv"
	"clean-arch/infra"
	"clean-arch/notification"
	"clean-arch/todo"
	"clean-arch/user"
	"clean-arch/worker"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	db := infra.NewMemoryDB()
	m := infra.NewFakeMailSender()
	n := notification.NewMailNotification("todo@mail.com", m)
	r := csv.NewCsvReader()
	uuc := user.NewUseCase()
	tuc := todo.NewUseCase(db, r, n)

	startWorker(tuc, uuc)
	startServer(tuc)
}

func startWorker(tuc *todo.UseCase, ucc *user.UseCase) {
	w := worker.NewWorker(tuc, ucc)
	if err := w.ImportFromFile("todos.csv"); err != nil {
		log.Fatalln(err)
	}

	if err := w.Notify(); err != nil {
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
