package app

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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
	if vkreq.Secret != rb.Config.VK.Secret {
		log.Println("Unknown service")
		return
	}
	if vkreq.Type == Confirm {
		b := []byte(rb.Config.VK.ConfirmationResponse)
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
		start := time.Now()
		defer func() {
			log.Printf(fmt.Sprintf("request executed for %s", time.Since(start).String()))
		}()
		if vkreq.IsCommand() {
			command, params, err := rb.ParseCommand(vkreq)
			if err != nil {
				log.Println(err)
				return
			}
			result, err := command(params...)
			if err != nil {
				result = err.Error()
			}
			user, err := rb.DB.GetUser(vkreq.Object.Message.FromID)
			if err != nil {
				u, err := rb.FindUser(vkreq.Object.Message.FromID)
				if err != nil {
					log.Println(err)
				}
				user = u
				if err == nil {
					err := rb.DB.SetUser(vkreq.Object.Message.FromID, user)
					if err != nil {
						log.Println(err)
					}
				}
			}
			result = user + "\n" + result

			err = vkreq.SendResult(result, rb.Generator, rb.Config)
		}
	}

}

func (rb *RollBot) Homepage(w http.ResponseWriter, req *http.Request) {
	tmpl, err := template.ParseFiles("templates/homepage.html")
	if err != nil {
		log.Println(err)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Println(err)
		return
	}
}

func (rb *RollBot) GetUsers(w http.ResponseWriter, req *http.Request) {
	tmpl, err := template.ParseFiles("templates/users.html")
	if err != nil {
		log.Println(err)
		return
	}
	users, err := rb.DB.GetUsers()
	if err != nil {
		log.Println(err)
		return
	}
	err = tmpl.Execute(w, users)
	if err != nil {
		log.Println(err)
		return
	}
}
