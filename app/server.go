package main

import (
	"fmt"
	"io"
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

	server := newServer("0.0.0.0:6379")
	err := server.Listen()
	if err != nil {
		os.Exit(1)
	}
	conn, _ := server.Accept()
	defer conn.Close()
}

func newServer(addr string) *Server {
	return &Server{
		Addr: addr,
	}
}

func (s *Server) Listen() error {
	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
	}
	defer l.Close()
	s.Listener = l
	return nil
}

func (s *Server) Accept() (net.Conn, error) {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			return nil, err
		}
		go s.HandleConnection(conn)
	}
}

func (s *Server) HandleConnection(conn net.Conn) {
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer) // Use the bufio.Reader to read from the connection
		if n == 0 { return }
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Println("Error reading from connection: ", err.Error())
		}
		buffer = buffer[:n]
		fmt.Printf("buffer: %s\n", buffer)
		_, err = conn.Write([]byte("PONG\r\n"))
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			return
		}
	}
}

