package sftp

import (
	"errors"
	"log"
	"net"
)

var FileNotFoundResponse = newResponse(ResponseFileNotFound, nil)
var ErrorInvalidFileName = errors.New("inavlid file name")

type Client struct {
	Conn net.Conn
}

func NewSFTPClient() (*Client, error) {
	return &Client{nil}, nil
}

func (c *Client) Connect(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	c.Conn = conn
	return nil
}

func (c *Client) Close() {
	c.Conn.Close()
}

func (c *Client) SendRequest(r Request) error {
	_, err := writeReqRes(c.Conn, r)
	return err
}

func (c *Client) ReadResponse() (Response, error) {
	return readRequestResponse(c.Conn)
}

func (c *Client) GetFile(fname string) error {
	return GetFile(fname, c.Conn)
}


func MakeRequest(c net.Conn, req Request) (Response, error) {
	if _, err := writeReqRes(c, req); err != nil {
		return nil, err
	}
	return readRequestResponse(c)
}

func GetFile(fname string, c net.Conn) error {

	req, err := NewRequest(RequestGetFile, NewGetFileHeader(fname))
	if err != nil {
		log.Fatal(err)
	}

	return WriteRequest(c, req)
}
