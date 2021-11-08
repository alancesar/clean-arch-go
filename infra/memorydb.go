package infra

import (
	"clean-arch/todo"
	"errors"
)

type MemoryDb struct {
	currentID int
	todos     []todo.Todo
}

func NewMemoryDB() *MemoryDb {
	return &MemoryDb{
		currentID: 0,
	}
}

func (mdb *MemoryDb) Insert(t todo.Todo) (int, error) {
	t.ID = mdb.getNextID()
	mdb.todos = append(mdb.todos, t)
	return t.ID, nil
}

func (mdb MemoryDb) FindByID(id int) (todo.Todo, error) {
	for _, curr := range mdb.todos {
		if curr.ID == id {
			return curr, nil
		}
	}

	return todo.Todo{}, errors.New("todo not found")
}

func (mdb *MemoryDb) Update(t todo.Todo) error {
	if err := checkIfIsValidForUpdate(t); err != nil {
		return err
	}

	for index, curr := range mdb.todos {
		if curr.ID == t.ID {
			mdb.todos[index] = t
			return nil
		}
	}

	return errors.New("todo not found")
}

func (mdb MemoryDb) FindByUserAndDone(user string, done bool) ([]todo.Todo, error) {
	var output []todo.Todo

	for _, curr := range mdb.todos {
		if curr.User == user && curr.Done == done {
			output = append(output, curr)
		}
	}

	return output, nil
}

func (mdb *MemoryDb) getNextID() int {
	mdb.currentID++
	return mdb.currentID
}

func checkIfIsValidForUpdate(t todo.Todo) error {
	if t.ID == 0 {
		return errors.New("id is required")
	}

	return nil
}
