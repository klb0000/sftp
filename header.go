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
