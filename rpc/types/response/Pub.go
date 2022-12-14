package response

import (
	"encoding/json"
	"fmt"
)

type Pub struct {
	Pub string `json:"pub"`
}

func (p *Pub) ToString() string {
	str := fmt.Sprintf("Pub='%s'", p.Pub)
	return str
}

func (p *Pub) ParseSelfFromJson(valueBytes []byte) {
	err := json.Unmarshal(valueBytes, p)
	if err != nil {
		p.Pub = ""
	}
}
