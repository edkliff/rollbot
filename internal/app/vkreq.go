package app

import (
	"fmt"
	"github.com/edkliff/rollbot/internal/config"
	"github.com/edkliff/rollbot/internal/generator"
	"io/ioutil"
	"net/http"
	"net/url"
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

func (vkr VKReq) SendResult(text string, gen *generator.Generator, c config.Config) error {
	fmt.Println("------ send result ------")
	params := make(map[string]string)
	params["user_id"] = fmt.Sprintf("%d", vkr.Object.Message.FromID)
	params["random_id"] =  fmt.Sprintf("%d", gen.Random(10000000,2147483646))
	if vkr.Object.Message.PeerID != 0 {
		delete(params, "user_id")
		params["peer_id"] = fmt.Sprintf("%d", vkr.Object.Message.PeerID)
	}
	params["message"] = text
	response, err := SendWithParams("messages.send", params, c)
	if err != nil {
		return err
	}
	fmt.Println("RESPONSE:", string(response))
	fmt.Println("------ result sended ------")
	return nil
}


func SendWithParams(method string, params map[string]string,  c config.Config) ([]byte, error) {
	params["access_token"] = c.VK.Token
	params["v"] = c.VK.APIVersion
	address, err := url.Parse(c.VK.VKServer + method)
	if err != nil {
		return nil, err
	}
	query := address.Query()
	for k, v := range params {
		query.Set(k, v)
	}
	address.RawQuery = query.Encode()
	fmt.Println("request", address.String())
	response, err := http.Get(address.String())
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (vkr *VKReq) IsCommand() bool  {
	s := strings.Replace(vkr.Object.Message.Text, "\"", "", -1)
	vkr.Object.Message.Text = s
	if len(s) > 0 && s[0] == '/' {
		return true
	}
	return false
}