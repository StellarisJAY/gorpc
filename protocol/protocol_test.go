package protocol

import (
	"bytes"
	"testing"
)

func TestHeader_MessageType(t *testing.T) {
	header := Header([10]byte{})
	header.SetMessageType(Response)
	if header.MessageType() != Response {
		t.Log("failed to set message type to Request")
		t.FailNow()
	}
	header.SetMessageType(Ping)
	if header.MessageType() != Ping {
		t.Log("failed to set message type to Ping")
		t.FailNow()
	}
	// 测试修改 SerializeType 是否对MessageType造成影响
	header.SetSerializeType(ProtobufSerialize)
	if header.MessageType() != Ping {
		t.Log("message type override by serial type")
		t.FailNow()
	}
}

func TestHeader_SerializeType(t *testing.T) {
	header := Header([10]byte{})
	header.SetSerializeType(JSONSerialize)
	if header.SerializeType() != JSONSerialize {
		t.Log("failed to set serialize type to JSON")
		t.FailNow()
	}
	header.SetSerializeType(ProtobufSerialize)
	if header.SerializeType() != ProtobufSerialize {
		t.Log("failed to set serialize type to Protobuf")
		t.Log(header[1])
		t.FailNow()
	}
	// 测试修改 MessageType 是否对serializeType造成影响
	header.SetMessageType(Ping)
	if header.SerializeType() != ProtobufSerialize {
		t.Log("serialize type override by message type")
		t.FailNow()
	}
}

func TestHeader_Seq(t *testing.T) {
	header := Header([10]byte{})
	header.SetSeq(1000)
	if header.Seq() != 1000 {
		t.Log("initial set seq failed")
		t.FailNow()
	}
	header.SetSeq(20)
	if header.Seq() != 20 {
		t.Log("reset seq failed")
		t.FailNow()
	}
}

func TestMessage_WriteTo_ReadFrom(t *testing.T) {
	message := messagePool.Get().(*Message)
	message.SetSeq(111)
	message.SetMessageType(Response)
	message.SetSerializeType(JSONSerialize)
	message.PutMeta("meta-1", "hello")
	message.PutMeta("meta-2", "world")
	message.ServiceName = "hello-service"
	message.ServiceMethod = "SayHello"
	defer messagePool.Put(message)

	buffer := &bytes.Buffer{}
	_, err := message.WriteTo(buffer)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	msg2 := messagePool.Get().(*Message)
	_, err = msg2.ReadFrom(buffer)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if msg2.Metadata == nil || msg2.Metadata["meta-1"] != "hello" || msg2.Metadata["meta-2"] != "world" {
		t.Log("metadata read write failed")
		t.FailNow()
	}
	if msg2.ServiceName != message.ServiceName || msg2.ServiceMethod != message.ServiceMethod {
		t.Log("service name or method read write failed")
		t.FailNow()
	}
	if msg2.Header.Seq() != 111 {
		t.Log("message header sequence number read write failed")
		t.FailNow()
	}
	if msg2.Header.MessageType() != Response || msg2.Header.SerializeType() != JSONSerialize {
		t.Log("header m_type or s_type read write failed")
		t.FailNow()
	}
}
