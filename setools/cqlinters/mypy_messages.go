package cqlinters

import (
	"fmt"
	"strings"
)

type MypyMessage struct {
	Severity string
	Message  string
	Filename string
	Line     int
	Column   int
}

func (msg MypyMessage) String() string {
	if msg.Line > 0 {
		return fmt.Sprint("`", msg.Filename, ":", msg.Line, ",", msg.Column, "` - ", strings.Title(msg.Severity), ": ", msg.Message)
	}
	return fmt.Sprint("`", msg.Filename, "` - ", strings.Title(msg.Severity), ": ", msg.Message)
}
