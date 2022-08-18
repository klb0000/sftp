package sftp

import (
	"encoding/json"
	"fmt"
)

const (
	RequestGetChunk = 129
	RequestGetFile  = 128
	RequestPut      = 160
	RequestDelete   = 192
	RequestQuery    = 224
)

var CodeToRequest = map[int]string{
	10: "Get",
	20: "put",
}


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

	Payload() []byte
}

// [version 1-byte][code 1-byte][HeaderLen 2-bytes][Header byte x-bytes]
// implementation of Request
type ftpRequest struct {
	version   uint8  // 1 byte
	code      uint8  // 1 byte
	HeaderLen uint16 // 2 byte

	//header is a json encoded data
	header  []byte
	payload []byte
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

func (r *ftpRequest) Payload() []byte {
	return r.payload
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
	return NewRequest(RequestGetChunk, h)

}

func NewRequest(code int, h Header) (Request, error) {

	HeaderByte, err := h.MarshalJson()
	if err != nil {
		return nil, err
	}

	return &ftpRequest{
		version:   1,
		code:      uint8(code),
		HeaderLen: uint16(len(HeaderByte)),
		header:    HeaderByte,
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
