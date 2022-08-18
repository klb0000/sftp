package sftp

import (
	"encoding/binary"
	"fmt"
	"testing"
)

var TestHeader = []Header{
	{
		"FileName":   "file.txt",
		"FileHash":   nil,
		"Compress":   false,
		"encrpytion": "AES",
	},
	{
		"FileName": nil,
		"FileHash": "2a6266cd228e2f88999c",
		"Compress": true,
	},
	{

		"FileName": nil,
		"FileHash": nil,
		"Compress": true,
	},
}

var GetRequest, _ = NewRequest(RequestGetFile, TestHeader[0])


func TestMarshalBinary(t *testing.T) {
	var tests = []Header{TestHeader[0], TestHeader[1]}
	for i := range tests {
		h := tests[i]
		r, err := NewGetRequest(h)
		if err != nil {
			t.Error("invalid test instance")
			continue
		}

		b, _ := r.MarshalBinary()
		Header, _ := h.MarshalJson()
		if b[1] != RequestGetChunk || len(Header) != int(binary.BigEndian.Uint16(b[2:5])) {
			t.Error("error in binaryMarshalling")

		}
		fmt.Println(r)
		fmt.Println(string(Header))
		fmt.Println(b)
	}
}

func TestValidRequestHeader(t *testing.T) {
	var tests = []struct {
		input    Header
		expected bool
	}{
		{TestHeader[0], true},
		{TestHeader[2], false},
	}
	for i := range tests {
		h := tests[i].input
		got := IsValidGetRequestHeader(h)
		if got != tests[i].expected {
			t.Errorf("\ninput: RequestHeader(%v)\nexpected:  %v\ngot: %v\n",
				tests[i].input, tests[i].expected, got)
		}
	}

}





func TestMain(t *testing.T) {
	req, _ := NewGetRequest(
		Header{
			"FileName":    "file.mov",
			"hash":        "212ab323c3ef23fs323j9",
			"Encryption":  nil,
			"Compression": nil,
		},
	)
	fmt.Println(req)

}
