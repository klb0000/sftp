package sftp

import (
	"testing"
)

func TestSetKey(t *testing.T) {
	h := NewGetFileHeader("file")

	//ContentSize key should have value type int
	if err := h.Set(ContentSize, "33"); err != ErrorInvalidValueType {
		t.Error("should return typeError")
	}
	//Compression key should have value type int
	if err := h.Set(Compression, 1); err != ErrorInvalidValueType {
		t.Error("should return typeError")
	}
	//FileName key should have value type int
	if err := h.Set(FileName, 33939); err != ErrorInvalidValueType {
		t.Error("should return typeError")
	}

	if err := h.Set(FileName, "some file.txt"); err != nil {
		t.Error("should throw typeError")
	}
	if err := h.Set(Compression, true); err != nil {
		t.Error("should throw typeError")
	}
	if err := h.Set(ContentSize, uint64(33)); err != nil {
		t.Error("should throw typeerror")
	}
}
