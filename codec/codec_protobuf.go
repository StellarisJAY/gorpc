package codec

import (
	"fmt"
	"google.golang.org/protobuf/proto"
)

type ProtobufCodec struct {
}

func (c ProtobufCodec) Decode(data []byte, v interface{}) error {
	if m, ok := v.(proto.Message); ok {
		return proto.Unmarshal(data, m)
	}
	return fmt.Errorf("%T is not an implementation of proto.Message", v)
}

func (c ProtobufCodec) Encode(v interface{}) ([]byte, error) {
	if m, ok := v.(proto.Message); ok {
		return proto.Marshal(m)
	}
	return nil, fmt.Errorf("%T is not an implementation of proto.Message", v)
}
