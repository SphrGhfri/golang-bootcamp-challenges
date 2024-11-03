package server

import (
	"QueraMQ/queue"
	"net"
)

type Topic struct {
	Name    string
	MQ      queue.IMessageQueue
	clients []net.Conn
	close   chan bool
}

func NewTopic(name string) *Topic {
	return &Topic{
		Name:    name,
		MQ:      queue.NewMessageQueue(),
		clients: make([]net.Conn, 0),
	}
}

func (t *Topic) AddClient(conn net.Conn) {
	t.clients = append(t.clients, conn)
}

func (t *Topic) Close() {
	t.close <- true
}
