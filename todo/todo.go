package todo

import (
	"errors"
	"time"
)

type Todo struct {
	ID      int
	Title   string
	User    string
	DueDate *time.Time
	Done    bool
}

var AlreadyDoneErr = errors.New("todo already done")

func NewTodo(title, user string, dueDate *time.Time) Todo {
	return Todo{
		Title:   title,
		User:    user,
		DueDate: dueDate,
	}
}

func (t *Todo) MarkAsDone() error {
	if t.Done {
		return AlreadyDoneErr
	}

	t.Done = true
	return nil
}

func (t *Todo) MarkAsPending() error {
	if !t.Done {
		return errors.New("todo already pending")
	}

	t.Done = false
	return nil
}
