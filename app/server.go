package main

import (
	"fmt"
	"net"
	"os"
)

type Server struct {
	Addr     string
	Listener net.Listener
}

func main() {
	fmt.Println("Logs from your program will appear here!")
	server := *&Server{
		Addr: "0.0.0.0:6379",
	}
	server.Start()
}

func (s *Server) Start() {
	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
	}
	defer l.Close()
	s.Listener = l
	conn, err := s.Listener.Accept()	
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	s.HandleConnection(conn)
}

func (s *Server) HandleConnection(conn net.Conn) {
	for {
		buffer := make([]byte, 128)
		_, err := conn.Read(buffer) // Use the bufio.Reader to read from the connection
		if err != nil {
			fmt.Println("Error reading from connection: ", err.Error())
			os.Exit(1)
		}
		conn.Write([]byte("+PONG\r\n"))
	}
}

