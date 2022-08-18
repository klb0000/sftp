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
	return writeReqRes(c.conn, r)
}

func (c *Client) ReadResponse() (Response, error) {
	return readRequestResponse(c.conn)
}

func (c *Client) RequestChunk(Filehash []byte, chunkNumber int) error {
	req, err := NewRequest(
		GetChunkRequest,
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
