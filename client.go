package sftp

import (
	"fmt"
	"net"
)

type Client struct {
	conn net.Conn
}

func NewSFTPClient() (*Client, error) {
	return &Client{nil}, nil
}

func (c *Client) Connect(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	fmt.Printf("connected to server %v", conn)
	c.conn = conn
	return nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) SendRequest(r Request) error {
	_, err := writeReqRes(c.conn, r)
	return err
}

func (c *Client) ReadResponse() (Response, error) {
	return readRequestResponse(c.conn)
}

func (c *Client) GetFile(fname string) error {
	return GetFile(fname, c.conn)
}

func (c *Client) RequestChunk(Filehash []byte, chunkNumber int) error {
	req, err := NewRequest(
		RequestGetChunk,
		Header{
			"Filehash":    string(Filehash),
			"ChunkNumber": chunkNumber,
		},
	)
	if err != nil {
		return err
	}
	return c.SendRequest(req)

}

func MakeRequest(c net.Conn, req Request) (Response, error) {
	if _, err := writeReqRes(c, req); err != nil {
		return nil, err
	}
	return readRequestResponse(c)
}
