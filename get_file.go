package sftp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type DownloadInfo struct {
	FileHash   string
	ChunkSize  int
	TotalChunk int
}

var FileNotFoundResponse, _ = NewResponse(ResponseFileNotFound, nil)

func GetFile(fname string, c net.Conn) error {
	req, err := NewRequest(
		RequestGetFile,
		Header{
			"FileName": fname,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	res, err := MakeRequest(c, req)
	if err != nil {
		log.Fatal(err)
	}
	if res.Code() >= ResponseError {
		fmt.Println(res)
		return errors.New("file not found")
	}
	fmt.Println(res)
	return nil
}

func HandleGetFileRequest(req Request, c net.Conn) error {
	reqHeader := toHeaderMap(req.Header())
	if reqHeader == nil {
		return errors.New("invalid Header")
	}
	fname, err := getFileName(reqHeader)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Open(fname)
	if err != nil {
		WriteResponse(c, FileNotFoundResponse)
		return err
	}
	io.Copy(c, f)
	f.Close()
	c.Close()
	return nil

}

func getFileName(h Header) (string, error) {
	if name, ok := h["FileName"].(string); ok {
		return name, nil
	}
	return "", errors.New("FileNameNotFound")

}

func toHeaderMap(data []byte) Header {
	var m Header
	if err := json.Unmarshal(data, &m); err != nil {
		return nil
	}
	return m

}
