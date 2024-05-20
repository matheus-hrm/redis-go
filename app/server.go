package main

import (
	"bufio" // Import the bufio package
	"fmt"
	"net"
	"os"
)

type Server struct {
	Addr     string
	Port     int
	Listener net.Listener
}

func main() {
	fmt.Println("Logs from your program will appear here!")

	server := newServer("0.0.0.0", 6379)
	err := server.Listen()
	if err != nil {
		os.Exit(1)
	}
	conn, _ := server.Accept()
	server.ReadLoop(conn)
}

func newServer(addr string, port int) *Server {
	return &Server{
		Addr: addr,
		Port: port,
	}
}

func (s *Server) Listen() error {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Addr, s.Port))
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
	}
	defer l.Close()
	s.Listener = l
	return nil
}

func (s *Server) Accept() (net.Conn, error) {
	conn, err := s.Listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		return nil, err
	}
	return conn, nil
}

func (s *Server) ReadLoop(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn) // Create a bufio.Reader to read from the connection
	for {
		buf := make([]byte, 1024)
		n, err := reader.Read(buf) // Use the bufio.Reader to read from the connection
		if err != nil {
			fmt.Println("Error reading from connection: ", err.Error())
			return
		}
		fmt.Println("Received ", string(buf[:n]))
		conn.Write([]byte("+PONG\r\n"))
	}
}

