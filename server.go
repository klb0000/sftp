package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"
)

var commandCode = map[byte]string{
	1: "Download",
	2: "Upload",
}

var funcs = map[byte]func(c net.Conn, arg string){
	1: sendFile,
	2: recieveFile,
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
		}
		go handleConn(conn)
	}
}

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		fmt.Fprintf(os.Stdin, "listening for command...\n")
		cmd, arg := listenForCommand(c)
		fmt.Fprintf(os.Stdout, "%s %s\n", commandCode[cmd], arg)
		time.Sleep(1 * time.Second)
		funcs[cmd](c, string(arg))
	}
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
	handleError(err)
	info := ""
	for _, file := range files {
		info += file.Name()
		info += " "
	}
	info += "\n"
	return []byte(info)

}
