package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

type Server struct {
	Addr     string
	Listener net.Listener
}

func main() {
	fmt.Println("Logs from your program will appear here!")
	server := Server{
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
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go s.HandleConnection(conn)
	}
}

func (s *Server) HandleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		buffer := make([]byte, 128)
		n, err := conn.Read(buffer) 
		if err != nil {
			fmt.Println("Error reading from connection: ", err.Error())
			break
		}
		buffer = buffer[:n]
		fmt.Println("Received: ", string(buffer))

		request, err := ParseCommand(buffer)
		if err != nil {
			fmt.Println("Error parsing request: ", err.Error())
			break
		}
		//_, err = conn.Write([]byte("+PONG\r\n"))
		
		err = WriteCommand(conn, request)
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			break
		}
	}
}

func WriteCommand(conn net.Conn, req Request) (error) {
	switch req.Command {
	case ECHO:
		_, err := io.WriteString(conn, fmt.Sprintf("+%s\r\n", req.Args[1]))
		if err != nil {
			fmt.Println("Error responding to ECHO", err.Error())
			os.Exit(1)
		}
	case PING: 
		_, err := io.WriteString(conn, "+PONG\r\n")
		if err != nil {
			fmt.Println("Error responding to PING", err.Error())
			os.Exit(1)
		}
	}
	return nil
}

