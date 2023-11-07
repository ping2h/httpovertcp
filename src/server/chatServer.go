package server

import (
	"fmt"
	"log"
	"net"
)

type chatServer struct {
	host string
	port string
}

func NewChatServer(host string, port string) *chatServer {
	return &chatServer{
		host: host,
		port: port,
	}
}

func (s *chatServer) Run() {
	ln, err := net.Listen("tcp", net.JoinHostPort(s.host, s.port))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("server is running at port:", s.port)
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print(err)
		}
		go s.handler(conn)
	}
}

func (s *chatServer) handler(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(string(buf[:n-1]))
	}

}
