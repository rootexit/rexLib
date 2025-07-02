package rexErrors

import (
	"fmt"
	"github.com/rootexit/rexLib/rexCodes"
)

type CodeMsg struct {
	Code int32
	Msg  string
}

func (c *CodeMsg) Error() string {
	return fmt.Sprintf("code: %d, msg: %s", c.Code, c.Msg)
}

// New creates a new CodeMsg.
func New(code int32, msg string) error {
	return &CodeMsg{Code: code, Msg: msg}
}

func Quick(code int32, Lang string) error {
	return &CodeMsg{Code: code, Msg: rexCodes.StatusText(code, Lang)}
}
