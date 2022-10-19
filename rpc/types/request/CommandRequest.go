package request

import "fmt"

type CommandRequest struct {
	Command string        `json:"command"`
	Args    []interface{} `json:"args"`
}

func NewCommandRequest(command string, args []any) *CommandRequest {
	return &CommandRequest{
		Command: command,
		Args:    args,
	}
}

// TODO: enable it to print out any kind of args
func (c *CommandRequest) ToString() string {
	str := fmt.Sprintf("CommandRequest{command='%s', args=%s}", c.Command, c.Args)
	return str
}
