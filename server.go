package fileserver

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
)

// type FileServer interface {
// 	Listen(network string, addr string) (net.Listener, error)
// }

type FileServer struct {
	net.Listener
}

func NewServer(addr string) (*FileServer, error) {
	listner, err := net.Listen("tcp", addr)
	return &FileServer{listner}, err
}

func (s *FileServer) Serve() {
	for {
		conn, err := s.Accept()
		if err != nil {
			logError("", err)
			continue
		}

		fmt.Fprintf(os.Stdout, "new connection %s\n", conn)
		go handleConn(conn)
	}
}

// func main() {
// 	listener, err := net.Listen("tcp", "localhost:8000")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	for {
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			log.Print(err)
// 		}
// 		go handleConn(conn)
// 	}

func logError(msg string, err error) {
	if len(msg) > 1 {
		fmt.Fprint(os.Stderr, msg)
	}
	fmt.Fprintf(os.Stderr, "%s\n", err)
}

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		req, err := ListenIncomingRequest(c)
		if err != nil {
			logError(" bad request", err)
			continue
		}

		if err = handleRequest(req, c); err != nil {
			logError(" unable to handle request", err)
		}

	}

	// for {
	// 	fmt.Fprintf(os.Stdout, "listening for command...\n")
	// 	cmd, arg := listenForCommand(c)
	// 	fmt.Fprintf(os.Stdout, "%s %s\n", commandCode[cmd], arg)
	// 	time.Sleep(1 * time.Second)
	// 	funcs[cmd](c, string(arg))
	// }
}

func ListenIncomingRequest(c net.Conn) (Request, error) {
	fmt.Fprint(os.Stdout, "listening to request ...\n")
	return readRequest(c)
}

// reads bytes( of arbitary length) from reader
// first two byte is the request code
// second two byte is the length of request header bytes
func readRequest(c io.Reader) (*ftpRequest, error) {
	// read (4 bytes) request type and headers len
	buf := make([]byte, 4)
	n, err := c.Read(buf)
	if n != 4 || err != nil {
		return nil, errors.New("error parsing request")
	}
	code, headerLen := binary.BigEndian.Uint16(buf[:2]), int(binary.BigEndian.Uint16(buf[2:]))

	// read headers byte
	headersBuf := make([]byte, headerLen)
	_, err = c.Read(headersBuf)
	if err != nil {
		return nil, errors.New("error reading header")
	}

	req := &ftpRequest{code, headersBuf, nil}
	if !req.isValidRequest() {
		return nil, errors.New("invalid request")
	}
	return req, nil
}

func handleRequest(req Request, c net.Conn) error {
	return nil
}

var commandCode = map[byte]string{
	1: "Download",
	2: "Upload",
}

var funcs = map[byte]func(c net.Conn, arg string){
	1: sendFile,
	2: recieveFile,
}

func sendFile(c net.Conn, fileName string) {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	io.Copy(c, f)
	f.Close()
}

func recieveFile(conn net.Conn, fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		conn.Close()
		log.Fatal(err)
	}
	io.Copy(file, conn)
}

func listenForCommand(c net.Conn) (byte, string) {
	cmd := make([]byte, 2)
	c.Read(cmd)
	arg := make([]byte, int(cmd[1]))
	c.Read(arg)
	return cmd[0], string(arg)
}

func listDir() []byte {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
	info := ""
	for _, file := range files {
		info += file.Name()
		info += " "
	}
	info += "\n"
	return []byte(info)

}
