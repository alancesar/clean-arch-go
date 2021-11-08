package csv

import (
	"clean-arch/date"
	"clean-arch/todo"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
)

const (
	titleIndex      = 0
	emailIndex      = 1
	dueDateIndex    = 2
	expectedColumns = 2
)

type Reader struct {
}

func NewCsvReader() *Reader {
	return &Reader{}
}

func (Reader) Read(reader io.Reader) ([]todo.Todo, error) {
	lines, err := csv.NewReader(reader).ReadAll()
	if err != nil {
		return nil, err
	}

	var output []todo.Todo
	for _, line := range lines {
		if t, err := parseLineToTodo(line); err != nil {
			fmt.Println(err)
		} else {
			output = append(output, t)
		}
	}

	return output, nil
}

func parseLineToTodo(line []string) (todo.Todo, error) {
	if len(line) == expectedColumns {
		if line[dueDateIndex] != "" {
			return parseTodoWithDueDate(line)
		}

		return parseTodoWithoutDueDate(line)
	}

	return todo.Todo{}, errors.New("invalid format")
}

func parseTodoWithDueDate(line []string) (todo.Todo, error) {
	parsedDate, err := date.ParseFromString(line[dueDateIndex])
	if err != nil {
		return todo.Todo{}, err
	}

	return todo.NewTodo(line[titleIndex], line[emailIndex], &parsedDate), nil
}

func parseTodoWithoutDueDate(line []string) (todo.Todo, error) {
	return todo.NewTodo(line[titleIndex], line[emailIndex], nil), nil
}
