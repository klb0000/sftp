package main

import (
	"fmt"
	"log"
	"mylib/sftp"
	"net"
	"os"
)


func main() {
	addr := os.Args[1]
	server, err := sftp.NewFileServer(addr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(server)
	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
		go handleConn(conn)
	}

}

func handleConn(c net.Conn) {
	req, err := sftp.ReadRequest(c)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(req)
	handleRequest(c, req)
}

func handleRequest(c net.Conn, r sftp.Request) {
	h, _ := sftp.ConvertToHeaderMap(r.Header())
	defer c.Close()
	sftp.HandleGetFile(c, h) 
}

func handleGetFile(c net.Conn, h sftp.HeaderMap) {
	defer c.Close()
	fileName := sftp.ReadFileName(h)
	if len(fileName) < 1 {
			log.Fatal("file name not found")
	}
	if err := sftp.WriteFile(fileName, c); err != nil {
			return
	}

}
