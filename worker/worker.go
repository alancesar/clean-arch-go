package worker

import (
	"clean-arch/user"
	"fmt"
	"io"
	"log"
	"os"
)

type TodoService interface {
	NotifyWhenDueDateArrives(user string) error
	BatchInsert(reader io.Reader) error
}

type UserService interface {
	GetAllUsers() ([]user.User, error)
}

type Worker struct {
	ts TodoService
	us UserService
}

func NewWorker(todoService TodoService, userService UserService) *Worker {
	return &Worker{
		ts: todoService,
		us: userService,
	}
}

func (w Worker) ImportFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		_ = file.Close()
	}()

	return w.ts.BatchInsert(file)
}

func (w Worker) Notify() error {
	users, err := w.us.GetAllUsers()
	if err != nil {
		return err
	}

	for _, u := range users {
		if err := w.ts.NotifyWhenDueDateArrives(u.Email); err != nil {
			fmt.Println(err)
		}
	}

	return nil
}
