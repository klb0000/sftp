package sftp

import (
	"encoding/json"
	"fmt"
)

const (
	ResponseDownloadInfo  = 50
	ResponseFileChunk    = 0
	ResponseQueryAnswer  = 63
	ResponseError        = 64
	ResponseFileNotFound = 65
	PutResponseCode      = 20
)

var CodeToResponse = map[int]string{
	10: "Get",
	20: "put",
}

type Response interface {
	//protocol version
	Version() int

	// Response code
	Code() int

	//json encoded byte
	Header() []byte

	// binary stream of Response
	MarshalBinary() ([]byte, error)
	Payload() []byte

	// text representation of Response
	String() string
}

// [version 1-byte][code 1-byte][headersLen 2-bytes][headers byte x-bytes]
// implementation of Response
type ftpResponse struct {
	version      uint8  // 1 byte
	ResponseCode uint8  // 1 byte
	headersLen   uint16 // 2 byte

	//header is a json encoded data
	header  []byte
	payload []byte
}

func (r *ftpResponse) Code() int {
	return int(r.ResponseCode)
}

func (r *ftpResponse) Version() int {
	return int(r.version)
}

func (r *ftpResponse) Header() []byte {
	return r.header
}

func (r *ftpResponse) MarshalBinary() ([]byte, error) {
	return MarshalBinary(r)

}

func (r *ftpResponse) Payload() []byte {
	return r.payload
}

func (r *ftpResponse) String() string {
	s := fmt.Sprintf("sftp %d Response %s\n", r.version, CodeToResponse[r.Code()])
	var h Header
	json.Unmarshal(r.header, &h)
	for k, v := range h {
		s += fmt.Sprintf("%v: %v\n",
			k, v)
	}
	return s
}

func NewResponse(status int, h Header) (Response, error) {

	headersByte, err := h.MarshalJson()
	if err != nil {
		return nil, err
	}

	return &ftpResponse{
		version:      1,
		ResponseCode: uint8(status),
		headersLen:   uint16(len(headersByte)),
		header:       headersByte,
		payload:      nil,
	}, nil

}

// func IsValidHeader(h Header) bool {
// 	fileName := h["FileName"]
// 	hash := h["FileHash"]
// 	return !(fileName == nil && hash == nil)

// }

// func (r *ftpResponse) Body() []byte {
// 	return r.body
// }
// func (r *ftpResponse) isValidResponse() bool {
// 	//TODO
// 	return true
// }
