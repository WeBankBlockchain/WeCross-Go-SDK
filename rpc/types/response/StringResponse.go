package response

type StringResponse string

func (s *StringResponse) ToString() string {
	return string(*s)
}

func (s *StringResponse) ParseSelfFromJson(valueBytes []byte) {
	str := string(valueBytes)
	*s = StringResponse(str)
}
