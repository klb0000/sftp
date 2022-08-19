package sftp

import (
	"bytes"
	"compress/gzip"
	"io"
	"testing"
)

func TestReadNBytes(t *testing.T) {
	_data := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	stream10 := bytes.NewBuffer(_data)
	data, _ := readNbytes(stream10, 10)
	if !bytes.Equal(data, _data) {
		t.Error("read byte is different from the source")
	}

	stream5 := bytes.NewBuffer(_data[:5])
	_, err := readNbytes(stream5, 10)
	if err == nil {
		t.Error("should throw error trying to read more number bytes from source")
	}

}

func TestWriteReqRes(t *testing.T) {
	r := GetRequest
	var w bytes.Buffer
	n, err := writeReqRes(&w, r)
	if err != nil {
		t.Error(err)
	}
	b, _ := r.MarshalBinary()
	if n != len(b) {
		t.Errorf("invalid write count")
	}

	if !bytes.Equal(w.Bytes(), b) {
		t.Errorf("invalid write")
	}
}

func TestReadCodeAndHeaderLen(t *testing.T) {
	stream := bytes.NewBuffer([]byte{1, ResponseError, 0, 0xff})
	c, l := readCodeAndHeaderLen(stream)
	if c != ResponseError {
		t.Errorf("wrong code expected %d\n", ResponseError)
	}
	if l != 0xff {
		t.Errorf("wrong HeaderLen expected %d got: %d\n", 0xff, l)

	}

	// should return error
	streamLen3 := bytes.NewBuffer([]byte{1, ResponseError, 0})
	c, l = readCodeAndHeaderLen(streamLen3)
	if c != -1 || l != -1 {
		t.Errorf("should return error(-1) on stream of size less than 4")
	}
}

// func TestWriteFile(t *testing.T) {
// 	path := "readWrite.go"
// 	WriteFile(path, os.Stdout)
// }

func TestCompressWrite(t *testing.T) {
	data := []byte("this is a uncompressed data")
	r := bytes.NewBuffer(data)
	var w bytes.Buffer
	gr := gzip.NewWriter(&w)
	io.Copy(gr, r)
	// fmt.Println(w.Bytes())
}
