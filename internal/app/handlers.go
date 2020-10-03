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
	fmt.Println("--------- start message ---------")
	err = json.Unmarshal(body, vkreq)
	if err != nil {
		log.Println(err)
	}
	if vkreq.Type == Confirm {
		b := []byte("61543fb6")
		_, err = w.Write(b)
		if err != nil {
			log.Println(err)
		}
		return
	}
	b := []byte("ok")
	_, err = w.Write(b)
	if err != nil {
		log.Println(err)
	}
	if vkreq.Type == MessageNew {
		if vkreq.IsCommand() {
			command, params, err := rb.ParseCommand(vkreq)
			if err!= nil {
				log.Println(err)
				return
			}
			result, err := command(params...)
			if err != nil {
				result = err.Error()
			}
			user, err := rb.DB.GetUser(vkreq.Object.Message.FromID)
			if err != nil {
				result = err.Error()
			} else {
				result = user + "\n" + result
			}
			err = vkreq.SendResult(result, rb.Generator, rb.Config)
		}
	}
	fmt.Println("--------- finish message ---------")

}
