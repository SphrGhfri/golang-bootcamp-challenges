package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net"
	"strconv"
	"strings"
)

type Server struct {
	port string
}

type Response struct {
	Result string `json:"result"`
	Error  string `json:"error"`
}

func NewServer(port string) *Server {
	return &Server{port: port}
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on port", s.port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	var req map[string]string

	if err := decoder.Decode(&req); err != nil {
		s.respondWithError(conn, "Invalid request format")
		return
	}

	action, ok := req["action"]
	if !ok || (action != "add" && action != "sub") {
		s.respondWithError(conn, "Invalid action")
		return
	}

	numbers, ok := req["numbers"]
	if !ok || numbers == "" {
		s.respondWithError(conn, "'numbers' parameter missing")
		return
	}

	numList, err := s.parseNumbers(numbers)
	if err != nil {
		s.respondWithError(conn, "Invalid number format")
		return
	}

	result, err := s.calculate(action, numList)
	if err != nil {
		s.respondWithError(conn, err.Error())
		return
	}

	s.respondWithResult(conn, result)
}

func (s *Server) parseNumbers(numbers string) ([]int64, error) {
	numStrs := strings.Split(numbers, ",")
	numList := make([]int64, len(numStrs))

	for i, numStr := range numStrs {
		num, err := strconv.ParseInt(strings.TrimSpace(numStr), 10, 64)
		if err != nil {
			return nil, err
		}
		numList[i] = num
	}
	return numList, nil
}

func (s *Server) calculate(action string, numbers []int64) (int64, error) {
	var result int64

	for i, num := range numbers {
		if i == 0 {
			result = num
			continue
		}

		if action == "add" {
			if willOverflow(result, num, "add") {
				return 0, fmt.Errorf("Overflow")
			}
			result += num
		} else if action == "sub" {
			if willOverflow(result, num, "sub") {
				return 0, fmt.Errorf("Overflow")
			}
			result -= num
		}
	}
	return result, nil
}

func willOverflow(a, b int64, op string) bool {
	if op == "add" {
		return (b > 0 && a > math.MaxInt64-b) || (b < 0 && a < math.MinInt64-b)
	} else if op == "sub" {
		return (b < 0 && a > math.MaxInt64+b) || (b > 0 && a < math.MinInt64+b)
	}
	return false
}

func (s *Server) respondWithResult(conn net.Conn, result int64) {
	response := Response{
		Result: fmt.Sprintf("The result of your query is: %d", result),
		Error:  "",
	}
	json.NewEncoder(conn).Encode(response)
}

func (s *Server) respondWithError(conn net.Conn, errMsg string) {
	response := Response{
		Result: "",
		Error:  errMsg,
	}
	json.NewEncoder(conn).Encode(response)
}

func main() {
	server := NewServer("4001")
	server.Start()
}
