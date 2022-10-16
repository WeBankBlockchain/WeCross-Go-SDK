package eles

import "fmt"

type ChainErrorMessage struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

func (cem *ChainErrorMessage) ToString() string {
	str := fmt.Sprintf("ChainErrorMessage{chain='%s', message='%s'}", cem.Path, cem.Message)
	return str
}
