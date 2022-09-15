package codec

import "testing"

type RegisterRequest struct {
	Account  string   `json:"account"`
	Password string   `json:"password"`
	Age      *int     `json:"age"`
	Ratio    *float64 `json:"ratio"`
}

func TestJSONCodec_Decode(t *testing.T) {
	data := []byte("{\"account\":\"jay\", \"password\":\"123456\", \"age\":0, \"ratio\":1.34}")
	req := new(RegisterRequest)
	err := JSONCodec{}.Decode(data, req)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if req.Age == nil || *req.Age != 0 {
		t.Log("json decode integer type failed")
		t.FailNow()
	}
	if req.Account != "jay" || req.Password != "123456" {
		t.Log("json decode string type failed")
		t.FailNow()
	}
	if req.Ratio == nil || *req.Ratio != 1.34 {
		t.Log("json decode float failed")
		t.FailNow()
	}
}

func TestJSONCodec_Encode(t *testing.T) {
	excepted := "{\"account\":\"jay\",\"password\":\"123456\",\"age\":0,\"ratio\":1.34}"
	req := &RegisterRequest{
		Account:  "jay",
		Password: "123456",
	}
	age, ratio := 0, 1.34
	req.Age = &age
	req.Ratio = &ratio
	data, err := JSONCodec{}.Encode(req)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if string(data) != excepted {
		t.Log("json encode failed")
		t.Log(string(data))
		t.FailNow()
	}
}

func TestProtobufCodec_Encode(t *testing.T) {
	req := &PBRegisterRequest{
		Age:      0,
		Ratio:    1.34,
		Account:  "jay",
		Password: "123456",
	}
	data, err := ProtobufCodec{}.Encode(req)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(data)
}

func BenchmarkJSONCodec_Decode(b *testing.B) {
	codec := JSONCodec{}
	data := []byte("{\"account\":\"jay\",\"password\":\"123456\",\"age\":0,\"ratio\":1.34}")
	req := new(RegisterRequest)
	for i := 0; i < b.N; i++ {
		err := codec.Decode(data, req)
		if err != nil {
			b.Error(err)
			b.FailNow()
		}
	}
}

func BenchmarkJSONCodec_Encode(b *testing.B) {
	req := &RegisterRequest{
		Account:  "jay",
		Password: "123456",
	}
	age, ratio := 0, 1.34
	req.Age = &age
	req.Ratio = &ratio
	codec := JSONCodec{}
	for i := 0; i < b.N; i++ {
		_, err := codec.Encode(req)
		if err != nil {
			b.Error(err)
			b.FailNow()
		}
	}
}

func BenchmarkProtobufCodec_Decode(b *testing.B) {
	data := []byte{10, 3, 106, 97, 121, 18, 6, 49, 50, 51, 52, 53, 54, 37, 31, 133, 171, 63}
	codec := ProtobufCodec{}
	req := new(PBRegisterRequest)
	for i := 0; i < b.N; i++ {
		err := codec.Decode(data, req)
		if err != nil {
			b.Error(err)
			b.FailNow()
		}
	}
}

func BenchmarkProtobufCodec_Encode(b *testing.B) {
	req := &PBRegisterRequest{
		Age:      0,
		Ratio:    1.34,
		Account:  "jay",
		Password: "123456",
	}
	codec := ProtobufCodec{}
	for i := 0; i < b.N; i++ {
		_, err := codec.Encode(req)
		if err != nil {
			b.Error(err)
			b.FailNow()
		}
	}
}
