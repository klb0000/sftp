package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
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

type Command interface {
	Code() int
	Arg() string
}

type command struct {
	code int
	arg  string
}

func (cmd *command) Code() int {
	return cmd.code
}

func (cmd *command) Arg() string {
	return cmd.arg
}

type FtpConnection struct {
	net.Conn
}

func newFtpConnection(serverAddr string) (*FtpConnection, error) {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		return nil, err
	}
	return &FtpConnection{conn}, nil
}

func (c *FtpConnection) Send(cmd Command) {
	argLen := len(cmd.Arg())
	commandHeader := []byte{byte(cmd.Code()), byte(argLen)}
	c.Write(commandHeader)
	argBytes := []byte(cmd.Arg())
	c.Write(argBytes)
}

func (conn *FtpConnection) Run() {
	defer conn.Close()
	for {
		fmt.Fprintf(os.Stdin, "> ")
		cmd := getCommand(os.Stdin)
		if isValidCommand(cmd) {
			conn.Send(cmd)
			switch cmd.Code() {
			case Download:
				getFile(conn, cmd.Arg())
			}
		}
	}
}

//get command (usually from stdin)
func getCommand(source io.Reader) Command {
	scanner := bufio.NewScanner(source)
	scanner.Scan()
	code, arg := splitCommand(scanner.Text())
	return &command{
		code: commandCode[code],
		arg:  arg,
	}

}

func isValidCommand(cmd Command) bool {
	switch cmd.Code() {
	case 1:
		return isValidDownloadCommand(cmd)
	}
	return false
}

func isValidDownloadCommand(cmd Command) bool {
	return len(cmd.Arg()) >= 1
}

// split command text in command and argument
func splitCommand(command string) (string, string) {
	cmd := strings.Split(command, " ")

	if len(cmd) < 2 {
		return cmd[0], ""
	}
	return cmd[0], cmd[1]
}

func main() {
	conn, err := newFtpConnection("localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	conn.Run()
}

func getFile(conn net.Conn, fileName string) (int, error) {
	file, err := os.Create(fileName)
	if err != nil {
		return -1, err
	}
	defer file.Close()
	return writeBytes(file, conn)
}

func writeBytes(dst io.Writer, src io.Reader) (int, error) {
	bufSize := 512
	bytesRecieved := 0
	buf := make([]byte, bufSize)

	for {
		nRead, err := src.Read(buf)
		bytesRecieved += nRead
		fmt.Printf("  bytes Recieved : %d\r", bytesRecieved)

		_, wErr := dst.Write(buf[:nRead])
		if wErr != nil {
			fmt.Fprintf(os.Stderr, "%s", wErr)
			return bytesRecieved, err
		}

		if nRead < bufSize || err != nil {
			fmt.Println()
			break
		}
	}

	return bytesRecieved, nil
}
