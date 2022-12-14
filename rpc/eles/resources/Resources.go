package resources

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Resources struct {
	ResourceDetails []*ResourceDetail `json:"resourceDetails"`
}

func (r *Resources) ToString() string {
	str := fmt.Sprintf("Resources{resourceDetails=%s}", ListResourceDetailToString(r.ResourceDetails))
	return str
}

func (r *Resources) ParseSelfFromJson(valueBytes []byte) {
	//r.ResourceDetails = make([]*ResourceDetail, 0)
	//m := make(map[string]interface{})
	//err := json.Unmarshal(valueBytes, &m)
	//if err != nil {
	//	return
	//}
	//details, ok := m["resourceDetails"].([]interface{})
	//if !ok {
	//	return
	//}
	//for _, obj := range details {
	//	objBytes, errr := json.Marshal(obj)
	//	if errr != nil {
	//		continue
	//	}
	//	detail := new(ResourceDetail)
	//	errr = json.Unmarshal(objBytes, detail)
	//	if errr != nil {
	//		continue
	//	}
	//	r.ResourceDetails = append(r.ResourceDetails, detail)
	//}
	err := json.Unmarshal(valueBytes, r)
	if err != nil {
		r.ResourceDetails = make([]*ResourceDetail, 0)
	}

}

func ListResourceDetailToString(resourcesDetails []*ResourceDetail) string {
	str := "["
	for i := 0; i < len(resourcesDetails); i++ {
		str += resourcesDetails[i].ToString()
		str += ", "
	}
	str = strings.TrimSuffix(str, ", ")
	str += "]"
	return str
}
