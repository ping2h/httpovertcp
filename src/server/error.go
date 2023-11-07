package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func error501(conn net.Conn) {
	log.Println("received a not implemented request")
	file, err := os.Open("src/server/501.html")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	fmt.Fprint(conn, "HTTP/1.1 501 Not Implemented\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", int(fileInfo.Size()))
	fmt.Fprint(conn, "Content-Type: text/html\r\n")
	fmt.Fprint(conn, "\r\n")
	io.Copy(conn, file)
}
func error404(conn net.Conn) {
	log.Println("received request that aquires non exist resource")
	file, err := os.Open("src/server/404.html")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	fmt.Fprint(conn, "HTTP/1.1 404 Not Found\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", int(fileInfo.Size()))
	fmt.Fprint(conn, "Content-Type: text/html\r\n")
	fmt.Fprint(conn, "\r\n")
	io.Copy(conn, file)
}

func error400(conn net.Conn) {
	log.Println("received a Bad Request")
	file, err := os.Open("src/server/400.html")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	fmt.Fprint(conn, "HTTP/1.1 501 Bad Request\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", int(fileInfo.Size()))
	fmt.Fprint(conn, "Content-Type: text/html\r\n")
	fmt.Fprint(conn, "\r\n")
	io.Copy(conn, file)
}
