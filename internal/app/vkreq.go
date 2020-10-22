package app

import (
	"strings"
)

type VKReq struct {
	Type   ReqType `json:"type"`
	Object struct {
		Message struct {
			Date                  int           `json:"date"`
			FromID                int           `json:"from_id"`
			ID                    int           `json:"id"`
			Out                   int           `json:"out"`
			PeerID                int           `json:"peer_id"`
			Text                  string        `json:"text"`
			ConversationMessageID int           `json:"conversation_message_id"`
			FwdMessages           []interface{} `json:"fwd_messages"`
			Important             bool          `json:"important"`
			RandomID              int           `json:"random_id"`
			Attachments           []interface{} `json:"attachments"`
			IsHidden              bool          `json:"is_hidden"`
		} `json:"message"`
		ClientInfo struct {
			ButtonActions  []string `json:"button_actions"`
			Keyboard       bool     `json:"keyboard"`
			InlineKeyboard bool     `json:"inline_keyboard"`
			Carousel       bool     `json:"carousel"`
			LangID         int      `json:"lang_id"`
		} `json:"client_info"`
	} `json:"object"`
	GroupID int    `json:"group_id"`
	EventID string `json:"event_id"`
	Secret  string `json:"secret"`
}

func (vkr *VKReq) RemoveQuotesAndCheckIsCommand() bool {
	s := strings.Replace(vkr.Object.Message.Text, "\"", "", -1)
	vkr.Object.Message.Text = s
	if len(s) > 0 && s[0] == '/' {
		return true
	}
	return false
}
