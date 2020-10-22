package app

import (
	"fmt"
	"github.com/edkliff/rollbot/internal/config"
	"io/ioutil"
	"net/http"
	"net/url"
)

func (rb RollBot) SendResult(vkr *VKReq, text string) error {

	params := make(map[string]string)
	params["user_id"] = fmt.Sprintf("%d", vkr.Object.Message.FromID)
	params["random_id"] = fmt.Sprintf("%d", rb.Generator.Random(10000000, 2147483646))
	if vkr.Object.Message.PeerID != 0 {
		delete(params, "user_id")
		params["peer_id"] = fmt.Sprintf("%d", vkr.Object.Message.PeerID)
	}
	params["message"] = text
	if vkr.Object.Message.ConversationMessageID != 0 {
		params["reply_to"] = fmt.Sprintf("%d", vkr.Object.Message.ID)
	}
	_, err := SendMessage("messages.send", params, rb.Config)
	if err != nil {
		return err
	}

	return nil
}

func SendMessage(method string, params map[string]string, c config.Config) ([]byte, error) {
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
