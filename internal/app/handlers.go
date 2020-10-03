package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func (rb *RollBot) VKHandle(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
	}
	vkreq := &VKReq{}
	err = json.Unmarshal(body, vkreq)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(body))
	fmt.Println(vkreq)
	if vkreq.Type == Confirm {
		b := []byte("61543fb6")
		_, err = w.Write(b)
		if err != nil {
			log.Println(err)
		}
		return
	}
	err = vkreq.SendResult("Бот на обслуживании", rb.Generator, rb.Config)
	b := []byte("ok")
	_, err = w.Write(b)
	if err != nil {
		log.Println(err)
	}
}
