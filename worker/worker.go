package worker

import (
	"io"
	"log"
	"os"
)

type Service interface {
	BatchInsert(reader io.Reader) error
}

type Worker struct {
	s Service
}

func NewWorker(service Service) *Worker {
	return &Worker{
		s: service,
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
	
	return w.s.BatchInsert(file)
}
