package app

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"text/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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
		if vkreq.RemoveQuotesAndCheckIsCommand() {
			command, err := rb.ParseCommand(vkreq)
			if err != nil {
				log.Println(err)
				return
			}
			result, err := command(vkreq)
			if err != nil {
				result = NewErrorResult(err)
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
			err = rb.DB.WriteTask(vkreq.Object.Message.Text,
				result.HTML(), result.Comment(),
				vkreq.Object.Message.FromID)
			if err != nil {
				log.Println(err)
			}
			err = rb.SendResult( vkreq, user + "\n" +result.String())
		}
	}

}

func (rb *RollBot) Homepage(w http.ResponseWriter, req *http.Request) {
	tmpl, err := template.ParseFiles( "templates/homepage.html.tmpl","templates/base.html.tmpl")
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
	tmpl, err := template.ParseFiles("templates/users.html.tmpl", "templates/base.html.tmpl")
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


func (rb *RollBot) GetHistory(w http.ResponseWriter, req *http.Request) {
	tmpl, err := template.ParseFiles("templates/logs.html.tmpl","templates/base.html.tmpl")
	if err != nil {
		log.Println(err)
		return
	}
	users, err := rb.DB.GetLogs(0)
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


func (rb *RollBot) GetUserHistory(w http.ResponseWriter, req *http.Request) {
	tmpl, err := template.ParseFiles("templates/logs.html.tmpl","templates/base.html.tmpl" )
	if err != nil {
		log.Println(err)
		return
	}
	userId :=  chi.URLParam(req, "userId")
	uid, err := strconv.Atoi(userId)
	if err != nil {
		log.Println(err)
		return
	}
	users, err := rb.DB.GetLogs(uid)
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
