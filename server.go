package sftp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
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
	Serve()

	Listen() (net.Conn, error)
	// Close closes the listener.
	// Any blocked Accept operations will be unblocked and return errors.
	Close() error
	// Addr returns the listener's network address.
	Addr() net.Addr
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

	// if err = handleRequest(req, c); err != nil {
	// 	logError(" unable to handle request", err)
	// }

}

func ReadRequest(c net.Conn) (Request, error) {
	// fmt.Fprint(os.Stdout, "listening to request ...\n")
	return readRequestResponse(c)
}

func readRequestResponse(c io.Reader) (RequestResponse, error) {

	buf := make([]byte, 4)
	if n, err := c.Read(buf); n != 4 || err != nil {
		return nil, err
	}
	//sftpVersion, requestCode and headerLen
	v, code, l := buf[0], buf[1], binary.BigEndian.Uint16(buf[2:])

	// read Header byte
	Header := make([]byte, l)
	if n, err := c.Read(Header); err != nil && n != int(l) {
		return nil, err
	}

	req := &ftpRequest{
		version:   v,
		code:      code,
		HeaderLen: l,
		header:    Header,
	}

	return req, nil
}

func WriteResponse(w io.Writer, r Response) error {
	return writeReqRes(w, r)
}

func writeReqRes(w io.Writer, r RequestResponse) error {
	data, err := r.MarshalBinary()
	if err != nil {
		return err
	}
	src := bytes.NewBuffer(data)
	if _, err := io.Copy(w, src); err != nil {
		return err
	}
	return nil
}

// // func handleRequest(req Request, c net.Conn) error {
// // 	return nil
// // }

// // // var commandCode = map[byte]string{
// // // 	1: "Download",
// // // 	2: "Upload",
// // // }

// // // var funcs = map[byte]func(c net.Conn, arg string){
// // // 	1: sendFile,
// // // 	2: recieveFile,
// // // }

// // // func sendFile(c net.Conn, fileName string) {
// // // 	f, err := os.Open(fileName)
// // // 	if err != nil {
// // // 		log.Fatal(err)
// // // 	}
// // // 	io.Copy(c, f)
// // // 	f.Close()
// // // }

// // // func recieveFile(conn net.Conn, fileName string) {
// // // 	file, err := os.Create(fileName)
// // // 	if err != nil {
// // // 		conn.Close()
// // // 		log.Fatal(err)
// // // 	}
// // // 	io.Copy(file, conn)
// // // }

// // // func listenForCommand(c net.Conn) (byte, string) {
// // // 	cmd := make([]byte, 2)
// // // 	c.Read(cmd)
// // // 	arg := make([]byte, int(cmd[1]))
// // // 	c.Read(arg)
// // // 	return cmd[0], string(arg)
// // // }

// // // func listDir() []byte {
// // // 	files, err := ioutil.ReadDir(".")
// // // 	if err != nil {
// // // 		fmt.Fprintf(os.Stderr, "%s\n", err)
// // // 	}
// // // 	info := ""
// // // 	for _, file := range files {
// // // 		info += file.Name()
// // // 		info += " "
// // // 	}
// // // 	info += "\n"
// // // 	return []byte(info)
