package sftp

import (
	"fmt"
	"log"
	"net"
	"os"
)

// type that statisfies both request and response
type RequestResponse interface {
	Request
	Response
}

type Server interface {
	net.Listener
}

type fileServer struct {
	net.Listener
}

func NewFileServer(addr string) (Server, error) {
	listner, err := net.Listen("tcp", addr)
	return &fileServer{listner}, err
}

func (s *fileServer) Listen() (net.Conn, error) {
	return s.Accept()
}

func (s *fileServer) Serve() {
	for {
		conn, err := s.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		fmt.Fprintf(os.Stdout, "new connection %v\n", conn)
		go handleConn(conn)
	}
}

func handleConn(c net.Conn) {
	defer c.Close()

	req, err := ReadRequest(c)
	if err != nil {
		// logError(" bad request", err)
		return
	}
	fmt.Println(req)

}

func HandleGetFile(c net.Conn, h HeaderMap) error {

	fileName := ReadFileName(h)
	if len(fileName) < 1 {
		WriteResponse(c, FileNotFoundResponse)
		return ErrorInvalidFileName
	}

	return WriteFile(fileName, c)

}

