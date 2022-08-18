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

	// if err = handleRequest(req, c); err != nil {
	// 	logError(" unable to handle request", err)
	// }

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
