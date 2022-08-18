package sftp

import (
	"bytes"
	"fmt"
	"testing"
)

// func TestTmp(t *testing.T) {
// 	Header := map[string]string{
// 		"file": "tmp.txt",
// 	}
// 	HeaderBytes, _ := json.Marshal(Header)

// 	r := ftpRequest{2, HeaderBytes, nil}
// 	fmt.Println(r.ToBytes())
// 	fmt.Println(r.Header())

// }

func TestReadRequest(t *testing.T) {
	r, _ := NewGetRequest(
		Header{
			"FileName":    "file.txt",
			"Compression": nil,
		},
	)
	data, _ := r.MarshalBinary()
	req, _ := readRequestResponse(bytes.NewReader(data))
	fmt.Println(req)

}

// func TestReadRequest(t *testing.T) {
// 	buf := []byte{0, 1, 0, 2, 3, 5}
// 	reader := bytes.NewReader(buf)
// 	req, err := readRequest(reader)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	if req.Code() != 1 {
// 		t.Error("wrong request code")
// 	}
// 	if len(req.HeaderBytes) != 2 {
// 		t.Error("wrong Header len")
// 	}

// }
