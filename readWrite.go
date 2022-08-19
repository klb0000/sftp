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

func ReadResponse(stream net.Conn) (Response, error) {
	c, hlen := readCodeAndHeaderLen(stream)
	if c < 0 || hlen < 0 {
		return nil, errors.New("error reading Request/Response code")
	}
	h, err := readHeader(stream, hlen)
	if err != nil {
		return nil, errors.New("error reading Request/Response code")
	}
	
	return newResponse(c, toHeaderMap(h)), nil
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

// func writeToStream(dstStream io.Writer, srcStream io.Reader) (int64, error) {
// 	return io.Copy(dstStream, srcStream)
// }

func WriteRequest(w io.Writer, r Request) error {
	_, err := writeReqRes(w, r)
	return err
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

func WriteFile(path string, dst io.Writer) error {
	f, err := os.Open(path)
	if err != nil {
		WriteResponse(dst, FileNotFoundResponse)
		return err
	}

	WriteResponse(dst, newResponse(
		ResponseGetOk,
		HeaderMap{
			"File_Name":   path,
			"Compression": false,
		},
	))

	_, err = io.Copy(dst, f)
	return err
}
