package request

type StringRequest string

func NewStringRequest(str string) *StringRequest {
	temp := StringRequest(str)
	return &temp
}

func (s *StringRequest) ToString() string {
	return string(*s)
}
