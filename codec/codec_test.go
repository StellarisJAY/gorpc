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

func TestProtobufCodec_Decode(t *testing.T) {

}

func TestProtobufCodec_Encode(t *testing.T) {

}

func BenchmarkJSONCodec_Decode(b *testing.B) {

}

func BenchmarkJSONCodec_Encode(b *testing.B) {

}
