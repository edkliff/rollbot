package app

import (
	"fmt"
	"github.com/edkliff/rollbot/internal/config"
	"github.com/edkliff/rollbot/internal/generator"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
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
	params := make(map[string]string)
	params["user_id"] = fmt.Sprintf("%d", vkr.Object.Message.FromID)
	params["random_id"] =  fmt.Sprintf("%d", gen.Random(10000000,2147483646))
	if vkr.Object.Message.PeerID != 0 {
		delete(params, "user_id")
		params["peer_id"] = fmt.Sprintf("%d", vkr.Object.Message.PeerID)
	}
	params["message"] = text
	if vkr.Object.Message.ConversationMessageID != 0 {
		params["reply_to"] = fmt.Sprintf("%d", vkr.Object.Message.ID)
	}
	_, err := SendWithParams("messages.send", params, c)
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

func (rb *RollBot) FindUser(userId int) (string, error ) {
		method := "users.get"
		params := make(map[string]string)
		params["user_ids"] = strconv.Itoa(userId)
		response, err := SendWithParams(method, params, rb.Config)
		if err != nil {
			return "", err
		}

		fmt.Println(string(response))
		// body := make(map[string]interface{})
		// err = json.Unmarshal(response, &body)
		// if err != nil {
		// 	return "", err
		// }
		// r := body["response"].([]interface{})
		// resp := r[0].(map[string]interface{})
		// name = resp["first_name"].(string)
	return "", nil
}