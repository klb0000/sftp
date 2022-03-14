package fileserver

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestToBytes(t *testing.T) {
	req := ftpRequest{
		code:         300,
		headersBytes: []byte{2, 3, 4, 5},
		body:         []byte{3, 4, 4},
	}

	reqBytes, err := req.ToBytes()
	if err != nil {
		t.Error(err)
		return
	}

	codeByte := make([]byte, 2)
	binary.BigEndian.PutUint16(codeByte, req.code)
	if !bytes.Equal(reqBytes[:2], codeByte) {
		t.Error(" req code not same")
	}

}
