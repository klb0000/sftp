package sftp

import (
	"encoding/json"
	"errors"
	"log"
)

//standard header keys
const (
	FileName    string = "File_Name"
	FileHash    string = "File_Hash"
	Compression string = "Compression"
	ContentSize string = "Content_Size"
	ChunkNumber string = "Chunk_Number"
)

var IsValidHeaderKey = map[string]bool{
	FileName: true, FileHash: true, Compression: true, ContentSize: true, ChunkNumber: true,
}

var (
	ErrorInvalidHeaderKey = errors.New("invalid header key")
	ErrorKeyNotSet        = errors.New("key not set")
	ErrorInvalidValueType = errors.New("invalid value type")
)

type HeaderMap map[string]interface{}

func (h HeaderMap) Get(key string) (interface{}, error) {
	if !IsValidHeaderKey[key] {
		return nil, ErrorInvalidHeaderKey
	}
	v, ok := h[key]
	if !ok {
		return nil, ErrorKeyNotSet
	}
	return v, nil
}

func (h HeaderMap) Set(key string, value interface{}) error {
	if !IsValidHeaderKey[key] {
		return ErrorInvalidHeaderKey
	}
	h[key] = value
	switch key {
	case FileName, FileHash:
		if _, ok := value.(string); !ok {
			return ErrorInvalidValueType
		}
	case ContentSize, ChunkNumber:
		if _, ok := value.(uint64); !ok {
			return ErrorInvalidValueType
		}
	case Compression:
		if _, ok := value.(bool); !ok {
			return ErrorInvalidValueType
		}
	}
	return nil
}

func (h HeaderMap) MarshalJson() ([]byte, error) {
	return json.Marshal(h)
}

func (h HeaderMap) String() string {
	j, err := h.MarshalJson()
	if err != nil {
		log.Fatal(err)
	}
	return string(j)
}

func ReadFileName(h HeaderMap) string {
	name, ok := h[FileName].(string)
	if !ok {
		return ""
	}
	return name
}

func ReadContentSize(h HeaderMap) uint64 {
	n, ok := h[ContentSize].(uint64)
	if !ok {
		return 0
	}
	return n
}

func ConvertToHeaderMap(data []byte) (HeaderMap, error) {
	var h HeaderMap
	err := json.Unmarshal(data, &h)
	return h, err
}

func toHeaderMap(data []byte) HeaderMap {
	var h HeaderMap
	if err := json.Unmarshal(data, &h); err != nil {
		panic(err)
	}
	return h
}

func NewGetFileHeader(fileName string) HeaderMap {
	return HeaderMap{
		FileName:    fileName,
		Compression: false,
		FileHash:    nil,
	}
}
