package server

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"os"
	"strconv"
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
	uri    map[string]*sync.RWMutex
}

var allowedContentTypes = map[string]string{
	"text/html":  "html",
	"text/plain": "txt",
	"image/gif":  "gif",
	"image/jpeg": "jpg",
	"image/png":  "png",
	"text/css":   "css",
}

func NewMux() *Mux {
	return &Mux{
		rwLock: sync.RWMutex{},
		uri:    map[string]*sync.RWMutex{},
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
	m.uri[uri] = &sync.RWMutex{}
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
	s.Mux.UpdateMux("/")
	s.Mux.UpdateMux("/upload")

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
	reader := bufio.NewReaderSize(conn, 10<<20)
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
			fmt.Print(line)
			header = append(header, line)
			line, err = reader.ReadString('\n')
		}
		if err == io.EOF {
			log.Println("Done, this connection is closed.")
			return
		} else {

			s.mux(conn, header, reader)

		}

		log.Println("Request processed. Waiting for the next request...")
	}
}

func (s *Server) mux(conn net.Conn, header []string, reader *bufio.Reader) {
	// request line
	m := strings.Fields(header[0])[0] // method
	u := strings.Fields(header[0])[1] // uri
	// 501
	if m != "GET" && m != "POST" {
		error501(conn)
	}

	if m == "GET" {

		s.GetPage(conn, u)
	}
	if m == "POST" && u == "/upload" {
		s.upload(conn, header, reader)
	}

}

func (s *Server) GetPage(conn net.Conn, uri string) {
	if uri == "/" {
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
	} else if uri == "/files" {
		if err := s.RefreshFilePage(); err != nil {
			log.Fatal(err)
		}
		file, err := os.Open("src/server/files.html")
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
	} else {
		uri = strings.Replace(uri, "/", "", -1)
		if s.Mux.ReadMux(uri) {
			// lock
			s.Mux.uri[uri].RLock()
			defer s.Mux.uri[uri].RUnlock()
			file, err := os.Open("src/server/upload/" + uri)
			if err != nil {
				log.Println(err)
				return
			}
			defer file.Close()
			contentType := getContentType(uri)
			fileInfo, _ := file.Stat()
			fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
			fmt.Fprintf(conn, "Content-Length: %d\r\n", int(fileInfo.Size()))
			fmt.Fprint(conn, "Content-Type: "+contentType+"\r\n")
			fmt.Fprint(conn, "\r\n")
			io.Copy(conn, file)

		} else {
			fmt.Println(s.Mux.uri)
			error404(conn)
		}
	}
}

// race here
func (s *Server) upload(conn net.Conn, header []string, reader *bufio.Reader) {
	var contentType, contentLength, conTentDis string
	for _, line := range header {
		key := strings.Fields(line)[0]
		if key == "Content-Length:" {
			contentLength = strings.Fields(line)[1]
		} else if key == "Content-Type:" {
			contentType = strings.Fields(line)[1]
		} else if key == "Content-Disposition:" {
			conTentDis = strings.Fields(line)[2]
		}
	}

	// fmt.Println(conTentDis, contentLength, contentType)

	if _, ok := allowedContentTypes[contentType]; !ok {
		log.Println("The file format is not supported")
		error400(conn)
		return
	}
	fileName := getFileName(conTentDis)
	s.Mux.UpdateMux(fileName)
	// lock
	s.Mux.uri[fileName].Lock()
	defer s.Mux.uri[fileName].Unlock()

	file, err := os.Create("/home/dellzp/tmp/dslab1/src/server/upload/" + fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()
	intlength := 0
	intlength, _ = strconv.Atoi(contentLength)
	buffer := make([]byte, intlength)
	// fmt.Println(reader.Buffered())
	reader.Read(buffer)
	if _, err := file.Write(buffer); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) RefreshFilePage() error {
	tmpl := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>List Page</title>
	</head>
	<body>
		<h1>List of Items:</h1>
		<ul>
			{{range .}}
			<li><<a href="http://localhost:8080/{{.}}">{{.}}</a></li>
			{{end}}
		</ul>
	</body>
	</html>`

	// Create or open the output HTML file
	file, err := os.Create("/home/dellzp/tmp/dslab1/src/server/files.html")
	if err != nil {
		return err
	}
	defer file.Close()

	// Parse the HTML template
	t, err := template.New("list").Parse(tmpl)
	if err != nil {
		return err
	}

	// Execute the template with the list data and write to the file
	list := []string{}
	// lock
	s.Mux.rwLock.RLock()
	defer s.Mux.rwLock.RUnlock()
	for k, _ := range s.Mux.uri {
		if k != "/" && k != "/upload" {
			k = strings.Replace(k, "/", "", -1)
			list = append(list, k)
		}
	}

	err = t.Execute(file, list)
	if err != nil {
		return err
	}

	return nil
}
func methodAndURI(method, uri string) string {
	return method + " " + uri
}

func getFileName(str string) string {
	parts := strings.Split(str, "=")
	return parts[1]
}

func getContentType(uri string) string {
	parts := strings.Split(uri, ".")
	suffix := parts[1]
	for k, v := range allowedContentTypes {
		if v == suffix {
			return k
		}
	}

	return ""
}
