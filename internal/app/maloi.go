package app

import (
	"fmt"
)

type Maloi struct {
	owner string
}

func (app *RollBot) RollMaloi(vk VKReq) (Resulter, error) {
	users, err := app.DB.GetUsers()
	if err != nil {
		return nil, err
	}
	r := app.Generator.Random(0, int64(len(users.Users)))
	if r >= int64(len(users.Users)) {
		r--
	}
	user := users.Users[r]
	h := Maloi{owner:user.Username}
	return  &h, nil
}


func (h *Maloi) VKString() string {
	s := fmt.Sprintf("Очко Малого уходит %s", h.owner)
	return s
}

func (h *Maloi) Comment() string  {
	s := fmt.Sprintf("Очко Малого уходит %s", h.owner)
	return s
}

func (h *Maloi) HTML() string {
	s := fmt.Sprintf("Очко Малого уходит %s", h.owner)
	return "<div>" + s + "</div>"
}