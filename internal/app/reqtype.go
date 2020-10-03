package app

import (
	"strings"
)

type ReqType uint8

const (
	UnknownReqType ReqType = iota
	Confirm
	MessageTyping
	MessageNew
)

var ReqTypeText = map[ReqType]string{
	Confirm:       "confirmation",
	MessageTyping: "message_typing_state",
	MessageNew:    "message_new",
}

func (rt *ReqType) UnmarshalJSON(data []byte) error {
	s := strings.Replace(string(data), "\"", "", -1)
	switch s {
	case "confirmation":
		z := Confirm
		*rt = z
	case "message_typing_state":
		z := MessageTyping
		*rt = z
	case "message_new":
		z := MessageNew
		*rt = z
	default:
		z := UnknownReqType
		*rt = z
	}
	return nil
}
