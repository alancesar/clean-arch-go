package csv

import (
	"clean-arch/date"
	"clean-arch/todo"
	"encoding/csv"
	"fmt"
	"io"
)

type CsvReader struct {
}

func NewCsvReader() *CsvReader {
	return &CsvReader{}
}

func (CsvReader) Read(reader io.Reader) ([]todo.Todo, error) {
	lines, err := csv.NewReader(reader).ReadAll()
	if err != nil {
		return nil, err
	}

	var output []todo.Todo
	for _, line := range lines {
		if parsedDate, err := date.ParseFromString(line[1]); err != nil {
			fmt.Println(err)
		} else {
			output = append(output, todo.NewTodo(line[0], parsedDate))
		}
	}

	return output, nil
}
