package app

import (
	"github.com/edkliff/rollbot/internal/config"
	"github.com/edkliff/rollbot/internal/generator"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type VKReq struct {
	Type    ReqType `json:"type"`
	GroupID int     `json:"group_id"`
	Object  struct {
		Date                  int64  `json:"date"`
		FromID                int64  `json:"from_id"`
		ID                    int64  `json:"id"`
		Out                   int64  `json:"out"`
		PeerID                int64  `json:"peer_id"`
		Text                  string `json:"text"`
		ConversationMessageID int64  `json:"conversation_message_id"`
		Important             bool   `json:"important"`
		RandomID              int64  `json:"random_id"`
		IsHidden              bool   `json:"is_hidden"`
	} `json:"object"`
}

func (vkr VKReq) SendResult(text string, gen *generator.Generator, c config.Config) error {
	method := "messages.send"
	params := make(map[string]string)
	if vkr.Object.PeerID > 2000000000 {
		params["chat_id"] = strconv.Itoa(int(vkr.Object.PeerID - 2000000000))
	} else {
		params["reply_to"] = strconv.Itoa(int(vkr.Object.ConversationMessageID))
	}
	params["peer_id"] = strconv.Itoa(int(vkr.Object.PeerID))
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
