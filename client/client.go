package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

// "fmt"

const (
	Download = 1
	Upload   = 2
)

var commandCode = map[string]int{
	"dl": 1,
	"ul": 2,
}

type Command struct {
	Code int
	Arg  string
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	handleConn(conn)
}

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		fmt.Fprintf(os.Stdin, "> ")
		cmd := getCommand(os.Stdin)
		time.Sleep(1 * time.Second)
		cmd.send(c)
		getFile(c, cmd.Arg)
	}
}

func getFile(c net.Conn, fileName string) {
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	bufSize := 512
	bytesRecieved := 0
	buf := make([]byte, bufSize)
	for {
		n, _ := c.Read(buf)
		n2, _ := f.Write(buf[:n])
		bytesRecieved += n2
		fmt.Printf("  bytes Recieved : %d\r", bytesRecieved)
		if n < bufSize {
			fmt.Println()
			break
		}

	}
	f.Close()

}

func (cmd *Command) send(c net.Conn) {
	argLen := len(cmd.Arg)
	cmdBytes := []byte{byte(cmd.Code), byte(argLen)}
	c.Write(cmdBytes)
	argBytes := []byte(cmd.Arg)
	c.Write(argBytes)
}

func getCommand(source io.Reader) Command {
	scanner := bufio.NewScanner(source)
	scanner.Scan()
	text := scanner.Text()
	cmds := strings.Split(text, " ")
	cmd, arg := cmds[0], cmds[1]
	return Command{
		Code: commandCode[cmd],
		Arg:  arg,
	}

}
