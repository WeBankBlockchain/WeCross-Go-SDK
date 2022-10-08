package types

import "fmt"

type Request struct {
	version  string
	data     Data
	callBack *CallBack
}

func (req *Request) ToString() string {
	str := fmt.Sprintf("Request{version='%s'", req.version)
	if req.data == nil {
		str += ""
	} else {
		str += ", data=" + req.data.ToString()
	}
	str += "}"
	return str
}
