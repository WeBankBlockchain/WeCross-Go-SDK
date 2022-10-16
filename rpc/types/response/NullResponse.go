package response

type NullResponse struct {
}

func (nr *NullResponse) ToString() string {
	return ""
}
