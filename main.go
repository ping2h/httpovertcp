package main

import (
	"flag"
	"fmt"

	"github.com/ping2h/tcpserver/src/client"
	"github.com/ping2h/tcpserver/src/server"
)

func main() {
	serv := server.NewServer("localhost", "8080", 10)
	// chatserver := server.NewChatServer("localhost", "8080")
	mode := flag.String("mode", "", "Specify 'server' or 'client' mode")
	flag.Parse()
	switch *mode {
	case "server":
		serv.Run()
	case "client":
		client.Client()
	default:
		fmt.Println("Usage: myprogram -mode=server|client")
		return
	}
}
