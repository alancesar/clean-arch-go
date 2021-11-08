package notification

import (
	"clean-arch/todo"
	"fmt"
)

type Mail struct {
	to     string
	from   string
	sender Sender
}

type Sender interface {
	SendMail(to, from, message string) error
}

func NewMailNotification(from string, sender Sender) *Mail {
	return &Mail{
		from:   from,
		sender: sender,
	}
}

func (m Mail) Notify(todo todo.Todo) error {
	message := fmt.Sprintf("your todo item \"%s\" expires today", todo.Title)
	return m.sender.SendMail(todo.User, m.from, message)
}
