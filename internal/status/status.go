package status

import (
	"fmt"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/codes"
)

// Status represents an RPC status code, message.  It is immutable
// and should be created with New, Newf.
type Status struct {
	Code    int32
	Message string
}

// New returns a Status representing c and msg.
func New(c codes.Code, msg string) *Status {
	return &Status{Code: int32(c), Message: msg}
}

// Newf returns New(c, fmt.Sprintf(format, a...)).
func Newf(c codes.Code, format string, a ...interface{}) *Status {
	return New(c, fmt.Sprintf(format, a...))
}

// Err returns an error representing c and msg.  If c is OK, returns nil.
func Err(c codes.Code, msg string) error {
	return New(c, msg).Err()
}

// Errorf returns Error(c, fmt.Sprintf(format, a...)).
func Errorf(c codes.Code, format string, a ...interface{}) error {
	return Err(c, fmt.Sprintf(format, a...))
}

// StatusCode returns the status code contained in s.
func (s *Status) StatusCode() codes.Code {
	if s == nil {
		return codes.Success
	}
	return codes.Code(s.Code)
}

// StatusMessage returns the message contained in s.
func (s *Status) StatusMessage() string {
	if s == nil {
		return ""
	}
	return s.Message
}

// Err returns an immutable error representing s; returns nil if s.StatusCode() is OK.
func (s *Status) Err() error {
	if s.StatusCode() == codes.Success {
		return nil
	}
	return &Error{s: s}
}

func (s *Status) String() string {
	return fmt.Sprintf("rpc error: code = %s desc = %s", s.StatusCode(), s.StatusMessage())
}

// Error wraps a pointer of a status proto. It implements error and Status,
// and a nil *Error should never be returned by this package.
type Error struct {
	s *Status
}

func (e *Error) Error() string {
	return e.s.String()
}

// RPCStatus returns the Status represented by se.
func (e *Error) RPCStatus() *Status {
	return e.s
}
