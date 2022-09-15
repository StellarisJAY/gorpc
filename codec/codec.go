package codec

import "fmt"

// Codec 消息体数据的编解码接口
type Codec interface {
	Decode([]byte, interface{}) error
	Encode(interface{}) ([]byte, error)
}

type RawCodec struct {
}

func (c RawCodec) Decode(data []byte, v interface{}) error {
	if b, ok := v.([]byte); ok {
		copy(b[0:], data[0:])
		return nil
	}
	return fmt.Errorf("%T is not a byte slice", v)
}

func (c RawCodec) Encode(v interface{}) ([]byte, error) {
	if b, ok := v.([]byte); ok {
		return b[:], nil
	}
	return nil, fmt.Errorf("%T is not a byte slice", v)
}
