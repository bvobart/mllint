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
}

func (msg MypyMessage) String() string {
	return fmt.Sprint("`", msg.Filename, ":", msg.Line, "` - ", strings.Title(msg.Severity), ": ", msg.Message)
}
