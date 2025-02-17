package common

type Messages []MessageItem
type MessageItem struct {
	Name   string
	Value  string
	IsCode bool
}

type Notify interface {
	Send(message Messages) error
}
