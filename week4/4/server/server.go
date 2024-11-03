package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
)

type Server struct {
	Addr   string
	topics map[string]*Topic
	ln     net.Listener
	mu     sync.Mutex
}

func NewServer(address string) *Server {
	return &Server{
		Addr:   address,
		topics: make(map[string]*Topic),
	}
}

func (s *Server) Run() error {
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	s.ln = ln
	log.Printf("Server started at %s", s.Addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		fmt.Println("Client connected:", conn.RemoteAddr())

		go s.handleConnection(conn)
	}
}

func (s *Server) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, topic := range s.topics {
		topic.Close()
	}
	connections := s.GetClientConnections()
	for _, connection := range connections {
		connection.Close()
	}
	s.ln.Close()
}

func (s *Server) GetTopic(topicName string) (*Topic, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	topic, exists := s.topics[topicName]
	if exists {
		return topic, true
	} else {
		newTopic := NewTopic(topicName)
		s.topics[topicName] = newTopic
		return newTopic, false
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	for {
		var request map[string]interface{}
		if err := decoder.Decode(&request); err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			} else if err.Error() == "EOF" {
				log.Println("Connection closed by client")
				return
			} else {
				log.Printf("Failed to decode request: %v", err)
				return
			}
		}

		action, ok := request["action"].(string)
		if !ok {
			s.sendError(encoder, "unknown action")
			continue
		}

		switch action {
		case "publish":
			s.handlePublish(request, encoder)
		case "subscribe":
			s.handleSubscribe(request, encoder, conn)
		case "unsubscribe":
			s.handleUnsubscribe(request, encoder, conn)
		case "shutdown":
			s.Stop()
		case "close_connection":
			s.ConnectionClose(encoder, conn)
		default:
			s.sendError(encoder, "unknown action")
		}
	}
}

func (s *Server) handlePublish(request map[string]interface{}, encoder *json.Encoder) {

	messageData, ok := request["message"].(map[string]interface{})
	if !ok {
		s.sendError(encoder, "message is required")
		return
	}

	topicName, topicOk := messageData["topic"].(string)
	content, contentOk := messageData["content"].(string)
	priority, priorityOk := messageData["priority"].(float64)

	if !topicOk {
		s.sendError(encoder, "topic is required")
		return
	}
	if !contentOk {
		s.sendError(encoder, "message content is required")
		return
	}
	if !priorityOk {
		s.sendError(encoder, "priority is required")
		return
	}

	topic, exists := s.GetTopic(topicName)
	if !exists {
		topic = NewTopic(topicName)
		s.mu.Lock()
		s.topics[topicName] = topic
		s.mu.Unlock()

	}

	topic.MQ.PushMessage(content, int(priority))

	response := map[string]interface{}{"status": "ok"}
	encoder.Encode(response)
}

func (s *Server) handleSubscribe(request map[string]interface{}, encoder *json.Encoder, conn net.Conn) {

	topicName, ok := request["topic"].(string)
	if !ok {
		s.sendError(encoder, "topic is required")
		return
	}

	s.mu.Lock()
	_, exists := s.topics[topicName]
	if !exists {
		topic := NewTopic(topicName)
		s.topics[topicName] = topic
	}
	s.mu.Unlock()

	s.topics[topicName].AddClient(conn)

	response := map[string]interface{}{"status": "ok"}
	encoder.Encode(response)
}

func (s *Server) handleUnsubscribe(request map[string]interface{}, encoder *json.Encoder, conn net.Conn) {
	topicName, ok := request["topic"].(string)
	if !ok {
		s.sendError(encoder, "topic is required")
	}

	s.mu.Lock()
	s.RemoveClient(encoder, topicName, conn)
	s.mu.Unlock()

	response := map[string]interface{}{"status": "ok"}
	encoder.Encode(response)
}

func (s *Server) sendError(encoder *json.Encoder, message string) {
	errorResponse := map[string]interface{}{"error": message}
	encoder.Encode(errorResponse)
}

func (s *Server) GetClientConnections() []net.Conn {
	s.mu.Lock()
	defer s.mu.Unlock()

	connections := make([]net.Conn, 0)
	for _, topic := range s.topics {
		connections = append(connections, topic.clients...)
	}
	return connections
}

func (s *Server) ConnectionClose(encoder *json.Encoder, conn net.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, t := range s.topics {
		for i := 0; i < len(t.clients); i++ {
			if conn == t.clients[i] {
				t.clients = append(t.clients[:i], t.clients[i+1:]...)
				break
			}
		}
	}
	conn.Close()

	response := map[string]interface{}{"status": "ok"}
	encoder.Encode(response)
}

func (s *Server) RemoveClient(encoder *json.Encoder, topicName string, conn net.Conn) {
	for _, t := range s.topics {
		if topicName == t.Name {
			for i := 0; i < len(t.clients); i++ {
				if t.clients[i] == conn {
					t.clients = append(t.clients[:i], t.clients[i+1:]...)
					break
				}
			}
		}
	}

	response := map[string]interface{}{"status": "ok"}
	encoder.Encode(response)
}
