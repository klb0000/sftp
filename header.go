package sftp

import (
	"encoding/json"
	"log"
)

type Header map[string]interface{}

func (h Header) Get(key string) (interface{}, bool) {
	v, ok := h[key]
	return v, ok
}

func (h Header) MarshalJson() ([]byte, error) {
	return json.Marshal(h)
}

func (h Header) String() string {
	j, err := h.MarshalJson()
	if err != nil {
		log.Fatal(err)
	}
	return string(j)
}

func ReadFileName(h Header) string {
	name, ok := h["File_Name"].(string)
	if !ok {
		return ""
	}
	return name
}

func ReadPayloadSize(h Header) uint64 {
	n, ok := h["Payload_Size"].(uint64)
	if !ok {
		return 0
	}
	return n
}

func ConvertToHeader(data []byte) (Header, error) {
	var h Header
	err := json.Unmarshal(data, &h)
	return h, err
}

