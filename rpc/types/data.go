package types

type Data interface {
	ToString() string
}

type ResponseData interface {
	Data
	ParseSelfFromJson(valueBytes []byte)
}
