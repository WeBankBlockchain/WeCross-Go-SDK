package response

import (
	"encoding/json"
	"fmt"
)

type Stubs struct {
	StubTypes []string `json:"stubTypes"`
}

func (s *Stubs) ParseSelfFromJson(valueBytes []byte) {
	err := json.Unmarshal(valueBytes, s)
	if err != nil {
		s.StubTypes = make([]string, 0)
	}
}

func (s *Stubs) ToString() string {
	str := fmt.Sprintf("Stubs{stubTypes=%v}", s.StubTypes)
	return str
}
