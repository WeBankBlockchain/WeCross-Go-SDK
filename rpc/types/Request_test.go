package types

import "testing"

type testDataType struct {
	A string
	B int
}

func (tdt *testDataType) ToString() string {
	return ""
}

func TestRequest_ToJson(t *testing.T) {
	req := NewRequest(&testDataType{
		A: "",
		B: 0,
	})
	if string(req.ToJson()) != "{\"version\":\"1\",\"data\":{\"A\":\"\",\"B\":0}}" {
		t.Fatalf("fail in ToJson")
	}
	t.Logf(string(req.ToJson()))
}
