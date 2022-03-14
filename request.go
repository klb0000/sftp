package fileserver

import (
	"encoding/binary"
	"encoding/json"
	"errors"
)

// var (
// 	InvalidRequestError = errors.New("invalid request")

// )

type Request interface {
	Code() int
	Headers() map[string]string
	Body() []byte
}

// implementation of Request
type ftpRequest struct {
	code uint16
	//headeBytes is a json encoded data
	headersBytes []byte
	body         []byte
}

// func NewRequest(code int, headers map[string]string, body []byte) Request {
// 	headersBytes := json.Marshal(headers)
// 	return &ftpRequest{code, }
// }

func (r *ftpRequest) Code() int {
	return int(r.code)
}

func (r *ftpRequest) Headers() map[string]string {
	var header map[string]string
	if err := json.Unmarshal(r.headersBytes, header); err != nil {
		return nil
	}
	return header
}

func (r *ftpRequest) Body() []byte {
	return r.body
}

func (r *ftpRequest) ToBytes() ([]byte, error) {

	reqBytes := make([]byte, 4+len(r.headersBytes)+len(r.body))
	binary.BigEndian.PutUint16(reqBytes[:2], r.code)
	binary.BigEndian.PutUint16(reqBytes[2:4], uint16(len(r.headersBytes)))

	//copy headers
	if n := copy(reqBytes[4:], r.headersBytes); n != len(r.headersBytes) {
		return nil, errors.New(" headers encode error")
	}

	//copy body
	if r.body != nil {
		if n := copy(reqBytes[4+len(r.headersBytes):], r.body); n != len(r.body) {
			return nil, errors.New(" body encode error")
		}
	}
	return reqBytes, nil
}

func (r *ftpRequest) isValidRequest() bool {
	//TODO
	return true
}
