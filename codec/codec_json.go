package codec

import (
	"bytes"
	"encoding/json"
)

type JSONCodec struct {
}

func (c JSONCodec) Decode(data []byte, res interface{}) error {
	decoder := json.NewDecoder(bytes.NewBuffer(data))
	decoder.UseNumber()
	return decoder.Decode(res)
}

func (c JSONCodec) Encode(res interface{}) ([]byte, error) {
	return json.Marshal(res)
}
