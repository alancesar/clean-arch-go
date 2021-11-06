package todo

import (
	"errors"
	"fmt"
	"io"
)

type Repository interface {
	FindByID(id int) (Todo, error)
	FindByDone(done bool) ([]Todo, error)
	Insert(Todo) (int, error)
	Update(Todo) error
}

type Reader interface {
	Read(reader io.Reader) ([]Todo, error)
}

type UseCase struct {
	repository Repository
	reader     Reader
}

func NewUseCase(repository Repository, reader Reader) *UseCase {
	return &UseCase{
		repository: repository,
		reader:     reader,
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

func (uc UseCase) MarkAsDoneByID(id int) (Todo, error) {
	t, err := uc.repository.FindByID(id)
	if err != nil {
		return t, err
	}

	if err := t.MarkAsDone(); err != nil {
		return t, err
	}

	err = uc.repository.Update(t)
	return t, err
}

func (uc UseCase) ListAllPending() ([]Todo, error) {
	return uc.repository.FindByDone(false)
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
