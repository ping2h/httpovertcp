package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Server struct {
	host    string
	port    string
	maxConn int
	Mux     *Mux
}

type Mux struct {
	rwLock sync.RWMutex
	uri    map[string]bool
}

func NewMux() *Mux {
	return &Mux{
		rwLock: sync.RWMutex{},
		uri:    map[string]bool{},
	}
}
func NewServer(host, port string, maxConn int) *Server {
	return &Server{
		host:    host,
		port:    port,
		maxConn: maxConn,
		Mux:     NewMux(),
	}
}

func (m *Mux) UpdateMux(uri string) {
	m.rwLock.Lock()
	defer m.rwLock.Unlock()
	m.uri[uri] = true
}

func (m *Mux) ReadMux(uri string) bool {
	m.rwLock.RLock()
	defer m.rwLock.RUnlock()
	_, ok := m.uri[uri]
	return ok

}

func (s *Server) Run() {
	maxConnChan := make(chan struct{}, s.maxConn)
	ln, err := net.Listen("tcp", net.JoinHostPort(s.host, s.port))
	if err != nil {
		log.Fatalf("Listen:%v", err)
	}
	// registeration of new resources
	s.Mux.UpdateMux(methodAndURI("GET", "/"))
	s.Mux.UpdateMux(methodAndURI("POST", "/upload"))

	log.Println("server listen at port:", s.port)
	for {

		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept:", err)
		}

		maxConnChan <- struct{}{}
		go s.handleConnection(conn, maxConnChan)
	}
}

func (s *Server) handleConnection(conn net.Conn, maxConnChan chan struct{}) {
	defer func(c chan struct{}) {
		<-c
	}(maxConnChan)
	defer conn.Close()
	log.Println("New tcp connection.")
	reader := bufio.NewReader(conn)
	for {
		// keep connection alive
		header := []string{}
		var line string
		var err error
		line, err = reader.ReadString('\n')
		for err == nil {
			if line == "\r\n" {
				break
			}
			log.Print(line)
			header = append(header, line)
			line, err = reader.ReadString('\n')
		}
		if err == io.EOF {
			log.Println("Done, this connection is closed.")
			return
		} else {

			s.mux(conn, header)

		}

		log.Println("Request processed. Waiting for the next request...")
	}
}

func (s *Server) mux(conn net.Conn, header []string) {
	// request line
	m := strings.Fields(header[0])[0] // method
	u := strings.Fields(header[0])[1] // uri
	// 501
	if m != "GET" && m != "POST" {
		error501(conn)
	}
	// 404
	if !s.Mux.ReadMux(methodAndURI(m, u)) {
		error404(conn)
	}

	if m == "GET" && u == "/" {
		index(conn)
	}
	if m == "POST" && u == "/upload" {
		upload(conn)
	}
	// if m == "GET" && u == "/contact" {
	// 	contact(conn)
	// }
	// if m == "GET" && u == "/apply" {
	// 	apply(conn)
	// }
	// if m == "POST" && u == "/apply" {
	// 	applyProcess(conn)
	// }
}

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

func index(conn net.Conn) {

	file, err := os.Open("src/server/index.html")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", int(fileInfo.Size()))
	fmt.Fprint(conn, "Content-Type: text/html\r\n")
	fmt.Fprint(conn, "\r\n")
	io.Copy(conn, file)
}

func upload(conn net.Conn) {
	// Convert the net.Conn to a bufio.Reader
	reader := bufio.NewReader(conn)

	allowedContentTypes := map[string]string{
		"text/html":  "html",
		"text/plain": "txt",
		"image/gif":  "gif",
		"image/jpeg": "jpg",
		"image/png":  "png",
		"text/css":   "css",
	}

	// Process each part in the multipart request
	for {
		part, err := mpReader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading part:", err)
			return
		}

		// Check if the content type is allowed
		_, ok := allowedContentTypes[part.Header.Get("Content-Type")]
		if !ok {
			fmt.Fprint(conn, "HTTP/1.1 400 Bad Request\r\n")
			fmt.Fprint(conn, "\r\n")
			fmt.Fprint(conn, "Invalid file type")
			return
		}

		// Create a file in the "upload" directory with a unique name
		outFile, err := os.Create(filepath.Join("src/server/upload", part.FileName()))
		if err != nil {
			fmt.Fprint(conn, "HTTP/1.1 500 Internal Server Error\r\n")
			fmt.Fprint(conn, "\r\n")
			fmt.Fprint(conn, "Failed to create the file")
			return
		}
		defer outFile.Close()

		// Copy the part's content to the newly created file
		_, err = io.Copy(outFile, part)
		if err != nil {
			fmt.Fprint(conn, "HTTP/1.1 500 Internal Server Error\r\n")
			fmt.Fprint(conn, "\r\n")
			fmt.Fprint(conn, "Failed to save the file")
			return
		}

		fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
		fmt.Fprint(conn, "\r\n")
		fmt.Fprintf(conn, "File uploaded: %s\n", part.FileName())
	}
}
func methodAndURI(method, uri string) string {
	return method + " " + uri
}
