package resource

import "fmt"

type ResourceDetail struct {
	Path       string                 `json:"path"`
	Distance   int                    `json:"distance"`
	StubType   string                 `json:"stubType"`
	Properties map[string]interface{} `json:"properties"`
	CheckSum   string                 `json:"checkSum"`
}

func (rd *ResourceDetail) ToString() string {
	str := fmt.Sprintf("ResourceDetail{path='%s', distance=%d, stubType='%s', properties=%v, checksum='%s'}", rd.Path, rd.Distance, rd.StubType, rd.Properties, rd.CheckSum)
	return str
}
