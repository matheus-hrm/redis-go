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

		request, err := ParseCommand(buffer)
		if err != nil {
			fmt.Println("Error parsing request: ", err.Error())
			_, err = conn.Write([]byte ("-ERR " + err.Error() + "\r\n"))
			if err != nil {
				fmt.Println("Error writing to connection: ", err.Error())
			}
			continue
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
	var response string 
	switch req.Command {
	case ECHO:
		if len(req.Args) > 0 {
		 	response = fmt.Sprintf("+%s\r\n", req.Args[0])
		} else {
			response = "+\r\n"
		}
	case PING: 
		response = "+PONG\r\n"
	}
	_, err := conn.Write([]byte(response))
	if err != nil {
		return err
	}

	return nil
}
