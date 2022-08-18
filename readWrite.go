package sftp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"os"
)

var (
	ErrorReadingHeader = errors.New("error reading header data")
	ErrorReadingStream = errors.New("error reading data stream")
)

func ReadRequest(c net.Conn) (Request, error) {
	// fmt.Fprint(os.Stdout, "listening to request ...\n")
	return readRequestResponse(c)
}

func readRequestResponse(stream io.Reader) (RequestResponse, error) {

	code, headerLen := readCodeAndHeaderLen(stream)
	if code < 0 || headerLen < 0 {
		return nil, errors.New("error reading Request/Response code")
	}
	header, err := readHeader(stream, headerLen)
	if err != nil {
		return nil, err
	}

	req := &ftpRequest{
		version:   1,
		code:      uint8(code),
		HeaderLen: uint16(headerLen),
		header:    header,
	}

	return req, nil
}

// read N number of bytes from stream
func readNbytes(stream io.Reader, N int) ([]byte, error) {
	buf := make([]byte, N)
	n, err := stream.Read(buf)
	if n != N {
		return nil, ErrorReadingStream
	}
	return buf, err
}

func readCodeAndHeaderLen(stream io.Reader) (int, int) {
	buf, err := readNbytes(stream, 4)
	if err != nil {
		return -1, -1
	}
	code, headerLen := buf[1], binary.BigEndian.Uint16(buf[2:])
	return int(code), int(headerLen)
}

func readHeader(stream io.Reader, headerLen int) ([]byte, error) {
	// read header bytes
	hbytes, err := readNbytes(stream, headerLen)
	if err != nil {
		return nil, err
	}
	return hbytes, nil
}

func WriteResponse(w io.Writer, r Response) error {
	_, err := writeReqRes(w, r)
	return err
}

func writeReqRes(w io.Writer, r RequestResponse) (int, error) {
	data, err := r.MarshalBinary()
	if err != nil {
		return 0, err
	}
	src := bytes.NewBuffer(data)
	n, err := io.Copy(w, src)
	if err != nil {
		return int(n), err
	}
	return int(n), nil
}


func WriteFile(path string, w io.Writer) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, f)
	return err
}

