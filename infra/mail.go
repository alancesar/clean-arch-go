package infra

import "fmt"

type FakeMailSender struct {
}

func NewFakeMailSender() *FakeMailSender {
	return &FakeMailSender{}
}

func (FakeMailSender) SendMail(to, from, message string) error {
	fmt.Printf("\n!!! N E W   M A I L !!!\n\nHello %s, there is a message to you from %s:\n  --> %s\n\n",
		to, from, message)

	return nil
}
