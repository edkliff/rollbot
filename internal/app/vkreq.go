package app

import (
	"fmt"
	"github.com/edkliff/rollbot/internal/config"
	"github.com/edkliff/rollbot/internal/generator"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
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
	if len(vkr.Object.Message.Text)>1 && vkr.Object.Message.Text[0] != byte('\\') &&vkr.Object.Message.Text[1] != byte('/') {
		return nil
	}
	method := "messages.send"
	params := make(map[string]string)
	if vkr.Object.Message.PeerID > 2000000000 {
		params["chat_id"] = strconv.Itoa(vkr.Object.Message.PeerID - 2000000000)
	} else if vkr.Object.Message.FromID > 0 {
		params["user_id"] = strconv.Itoa(vkr.Object.Message.PeerID)
	} else {
		params["reply_to"] = strconv.Itoa(vkr.Object.Message.ConversationMessageID)
	}
	params["peer_id"] = strconv.Itoa(vkr.Object.Message.PeerID)
	params["random_id"] = strconv.Itoa(int(gen.Random(9000000, 9999999)))
	params["message"] = text


	_, err := SendWithParams(method, params, c)
	if err != nil {
		return err
	}

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
	fmt.Println("------------------")
	fmt.Println(address.RawQuery)
	response, err := http.Get(address.String())
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(content))
	fmt.Println("------------------")
	return content, nil
}
