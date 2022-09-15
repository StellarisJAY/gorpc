package protocol

import "testing"

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
