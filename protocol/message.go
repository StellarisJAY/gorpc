package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

// Header 协议首部，格式如下：
// magic | m_type | serial_type | seq |
// magic: 1 byte
// messageType: 4 bits
// serializeType: 4 bits
// sequenceNumber: 8 bytes
type Header [10]byte
type MessageType byte
type SerializeType byte

const (
	Request  MessageType = iota // Request 客户端请求消息类型
	Oneway                      // Oneway 单向消息，服务端不做回复
	Ping                        // Ping 心跳包
	Response                    // Response 正常返回消息类型
	Pong                        // Pong 心跳返回
	Error                       // Error 错误返回消息类型
)

const (
	Raw               SerializeType = iota // Raw 无序列化
	JSONSerialize                          // JSONSerialize 消息体JSON序列化
	ProtobufSerialize                      // ProtobufSerialize protocol buffer 序列化
)

const magicNumber = 0xfe

type Message struct {
	*Header
	Metadata      map[string]string // Metadata RPC消息元数据，记录trace等信息
	ServiceName   string            // ServiceName RPC请求的服务名称
	ServiceMethod string            // ServiceMethod RPC请求的方法名称
	Data          []byte            // Data 序列化后的数据
	buf           []byte
}

var (
	ErrMetadataKeyValueMissing = errors.New("incomplete key value pair in metadata")
)

func (h *Header) CheckMagicNumber() bool {
	return h[0] == magicNumber
}

func (h *Header) MessageType() MessageType {
	// 首部的第二个字节的高4位作为消息类型
	return MessageType(h[1]) >> 4
}

func (h *Header) SerializeType() SerializeType {
	// 首部的第二个字节的低4位作为序列化类型
	return SerializeType(h[1] & 0x0f)
}

func (h *Header) Seq() uint64 {
	return binary.BigEndian.Uint64(h[2:])
}

func (h *Header) SetSeq(seq uint64) {
	binary.BigEndian.PutUint64(h[2:], seq)
}

func (h *Header) SetMessageType(messageType MessageType) {
	h[1] = h[1]&0x0f | (byte(messageType) << 4)
}

func (h *Header) SetSerializeType(serializeType SerializeType) {
	h[1] = h[1]&0xf0 | byte(serializeType)
}

func (m *Message) WriteTo(writer io.Writer) (int64, error) {
	n, err := writer.Write((*m.Header)[:])
	if err != nil {
		return int64(n), err
	}
	meta := &bytes.Buffer{}
	encodeMetadata(m.Metadata, meta)
	snl, sml := uint32(len(m.ServiceName)), uint32(len(m.ServiceMethod))
	dl, ml := uint32(len(m.Data)), uint32(meta.Len())
	// 消息体部分的总长度
	totalLength := (4 + snl) + (4 + sml) + (4 + ml) + (4 + dl)
	lengthBuf := make([]byte, 4)
	// 写入总长度
	binary.BigEndian.PutUint32(lengthBuf, totalLength)
	n, err = writer.Write(lengthBuf)
	if err != nil {
		return int64(n), err
	}
	// 写入 ServiceName
	if n, err := writeString(writer, m.ServiceName, snl, lengthBuf); err != nil {
		return n, err
	}
	// 写入 ServiceMethod
	if n, err := writeString(writer, m.ServiceMethod, sml, lengthBuf); err != nil {
		return n, err
	}
	// 写入 Metadata
	if n, err := writeBytes(writer, meta.Bytes(), ml, lengthBuf); err != nil {
		return n, err
	}
	// 写入 Data
	if n, err := writeBytes(writer, m.Data, dl, lengthBuf); err != nil {
		return n, err
	}
	return int64(10 + totalLength), nil
}

func (m *Message) ReadFrom(reader io.Reader) (int64, error) {
	n, err := io.ReadFull(reader, (*m.Header)[:])
	if err != nil {
		return int64(n), err
	}
	lengthBuf := make([]byte, 4)
	n, err = io.ReadFull(reader, lengthBuf)
	if err != nil {
		return int64(n), err
	}
	totalLength := binary.BigEndian.Uint32(lengthBuf)
	var buffer []byte
	// 如果容量足够的话，重用buffer
	if cap(m.buf) >= int(totalLength) {
		buffer = m.buf[:totalLength]
	} else {
		buffer = make([]byte, totalLength)
		m.buf = buffer
	}
	n, err = io.ReadFull(reader, buffer)
	if err != nil {
		return int64(n), err
	}
	var index uint32 = 0
	snl := binary.BigEndian.Uint32(buffer[index : index+4])
	index += 4
	m.ServiceName = string(buffer[index : index+snl])
	index += snl
	sml := binary.BigEndian.Uint32(buffer[index : index+4])
	index += 4
	m.ServiceMethod = string(buffer[index : index+sml])
	index += sml
	ml := binary.BigEndian.Uint32(buffer[index : index+4])
	index += 4
	m.Metadata, err = decodeMetadata(buffer[index : index+ml])
	if err != nil {
		return int64(10 + index), err
	}
	index += ml
	dl := binary.BigEndian.Uint32(buffer[index : index+4])
	index += 4
	m.Data = buffer[index : index+dl]
	return int64(10 + totalLength), nil
}

func encodeMetadata(metadata map[string]string, buffer *bytes.Buffer) {
	length := make([]byte, 4)
	for key, value := range metadata {
		binary.BigEndian.PutUint32(length, uint32(len(key)))
		buffer.Write(length)
		buffer.Write([]byte(key))
		binary.BigEndian.PutUint32(length, uint32(len(value)))
		buffer.Write(length)
		buffer.Write([]byte(value))
	}
}

func decodeMetadata(buffer []byte) (map[string]string, error) {
	metadata := make(map[string]string)
	var i uint32 = 0
	n := uint32(len(buffer))
	for i < n {
		kl := binary.BigEndian.Uint32(buffer[i : i+4])
		i += 4
		if i+kl > n {
			return metadata, ErrMetadataKeyValueMissing
		}
		key := string(buffer[i : i+kl])
		i += kl
		vl := binary.BigEndian.Uint32(buffer[i : i+4])
		i += 4
		if i+vl > n {
			return metadata, ErrMetadataKeyValueMissing
		}
		value := string(buffer[i : i+vl])
		i += vl
		metadata[key] = value
	}
	return metadata, nil
}

func writeString(writer io.Writer, s string, l uint32, buf []byte) (int64, error) {
	return writeBytes(writer, []byte(s), l, buf)
}

func writeBytes(writer io.Writer, b []byte, l uint32, buf []byte) (int64, error) {
	// 写入serviceName
	binary.BigEndian.PutUint32(buf, l)
	n, err := writer.Write(buf)
	if err != nil {
		return int64(n), err
	}
	n, err = writer.Write(b)
	return int64(n), err
}
