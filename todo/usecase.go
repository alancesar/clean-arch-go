package todo

import (
	"clean-arch/date"
	"errors"
	"fmt"
	"io"
	"time"
)

var ForbiddenErr = errors.New("forbidden")

type Repository interface {
	FindByID(id int) (Todo, error)
	FindByUserAndDone(user string, done bool) ([]Todo, error)
	Insert(Todo) (int, error)
	Update(Todo) error
}

type Reader interface {
	Read(reader io.Reader) ([]Todo, error)
}

type Notification interface {
	Notify(todo Todo) error
}

type UseCase struct {
	repository   Repository
	reader       Reader
	notification Notification
}

func NewUseCase(repository Repository, reader Reader, notification Notification) *UseCase {
	return &UseCase{
		repository:   repository,
		reader:       reader,
		notification: notification,
	}
}

func (uc UseCase) CreateTodo(todo *Todo) error {
	if todo == nil {
		return errors.New("todo cannot be null")
	}

	createdID, err := uc.repository.Insert(*todo)
	if err != nil {
		return err
	}

	todo.ID = createdID
	return nil
}

func (uc UseCase) MarkAsDoneByUserAndID(user string, id int) (Todo, error) {
	t, err := uc.repository.FindByID(id)
	if err != nil {
		return t, err
	} else if t.User != user {
		return Todo{}, ForbiddenErr
	}

	if err := t.MarkAsDone(); err != nil {
		return t, err
	}

	err = uc.repository.Update(t)
	return t, err
}

func (uc UseCase) ListAllPending(user string) ([]Todo, error) {
	return uc.repository.FindByUserAndDone(user, false)
}

func (uc UseCase) NotifyWhenDueDateArrives(user string) error {
	todos, err := uc.repository.FindByUserAndDone(user, false)
	if err != nil {
		return err
	}

	now := time.Now()
	today := date.ParseToString(now)

	for _, todo := range todos {
		if isDueDate(today, todo.DueDate) {
			uc.notify(todo)
		}
	}

	return nil
}

func (uc UseCase) BatchInsert(reader io.Reader) error {
	todos, err := uc.reader.Read(reader)
	if err != nil {
		return err
	}

	for _, t := range todos {
		if _, err := uc.repository.Insert(t); err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

func (uc UseCase) notify(todo Todo) {
	if err := uc.notification.Notify(todo); err != nil {
		fmt.Println(err)
	}
}

func isDueDate(target string, dueDate *time.Time) bool {
	if dueDate == nil {
		return false
	}

	return date.ParseToString(*dueDate) == target
}
