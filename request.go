package sftp

import (
	"encoding/json"
	"fmt"
)

const (
	GetRequest      = 10
	GetChunkRequest = 11
	PutRequestCode  = 20
)

var CodeToRequest = map[int]string{
	10: "Get",
	20: "put",
}

// 1. get/download
// 2. upload/put
// 3. list

// get protocol
// 1. client send request get
// 2. server send response (ok) or (error: file not found)
//  	if error: process stop
//  	if ok: ok repsonse has metadata:
// .    	Totalsize
// chunk size
// number of chunk
// file hash
// .
// 3.1  client send request for first chunk
// 3. 2: server send first chunk
// 3. 3 : client send request for second chunk

type Request interface {
	//protocol version
	Version() int

	// request code
	Code() int

	//json encoded byte
	Header() []byte

	// binary stream of request
	MarshalBinary() ([]byte, error)

	// text representation of request
	String() string
}

// [version 1-byte][code 1-byte][HeaderLen 2-bytes][Header byte x-bytes]
// implementation of Request
type ftpRequest struct {
	version    uint8  // 1 byte
	code       uint8  // 1 byte
	HeaderLen uint16 // 2 byte

	//header is a json encoded data
	header []byte
}

func (r *ftpRequest) Code() int {
	return int(r.code)
}

func (r *ftpRequest) Version() int {
	return int(r.version)
}

func (r *ftpRequest) Header() []byte {
	return r.header
}

func (r *ftpRequest) MarshalBinary() ([]byte, error) {
	return MarshalBinary(r)
}

func (r *ftpRequest) String() string {
	s := fmt.Sprintf("sftp %d %s", r.version, CodeToRequest[r.Code()])
	var h Header
	json.Unmarshal(r.header, &h)
	for k, v := range h {
		s += fmt.Sprintf("\n%v: %v",
			k, v)
	}
	return s
}

func NewGetRequest(h Header) (Request, error) {
	return NewRequest(GetRequest, h)

}

func NewRequest(code int, h Header) (Request, error) {

	HeaderByte, err := h.MarshalJson()
	if err != nil {
		return nil, err
	}

	return &ftpRequest{
		version:    1,
		code:       uint8(code),
		HeaderLen: uint16(len(HeaderByte)),
		header:     HeaderByte,
	}, nil
}

func IsValidGetRequestHeader(h Header) bool {
	fileName := h["FileName"]
	hash := h["FileHash"]
	return !(fileName == nil && hash == nil)

}

// func (r *ftpRequest) Body() []byte {
// 	return r.body
// }
// func (r *ftpRequest) isValidRequest() bool {
// 	//TODO
// 	return true
// }
