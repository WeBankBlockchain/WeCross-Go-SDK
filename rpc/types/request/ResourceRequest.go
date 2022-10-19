package request

import "fmt"

type ResourceRequest struct {
	IgnoreRemote bool `json:"ignoreRemote"`
}

func NewResourceRequest(ignoreRemot bool) *ResourceRequest {
	return &ResourceRequest{IgnoreRemote: ignoreRemot}
}

func (rr *ResourceRequest) ToString() string {
	str := fmt.Sprintf("ResourceRequest{ignoreRemote=%t}", rr.IgnoreRemote)
	return str
}
