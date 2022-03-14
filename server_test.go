package fileserver

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
)

func TestTmp(t *testing.T) {
	buf := make([]byte, 4)
	fmt.Println(buf)
	binary.BigEndian.PutUint16(buf[:2], 258)
	fmt.Println(buf)

}

func TestNewServer(t *testing.T) {
	server, err := NewServer("localhost:8000")
	if err != nil {
		t.Error(err)
	}
	server.Close()
	// server, err = NewServer("192.168.3.6:8000")
	// if err != nil {
	// 	t.Error(err)
	// }
	// server.Close()
}

func TestReadRequest(t *testing.T) {
	buf := []byte{0, 1, 0, 2, 3, 5}
	reader := bytes.NewReader(buf)
	req, err := readRequest(reader)
	if err != nil {
		t.Error(err)
		return
	}
	if req.Code() != 1 {
		t.Error("wrong request code")
	}
	if len(req.headersBytes) != 2 {
		t.Error("wrong headers len")
	}

}
