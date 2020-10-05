package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type UserResponse struct {
	Response []struct {
		ID              int    `json:"id"`
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		IsClosed        bool   `json:"is_closed"`
		CanAccessClosed bool   `json:"can_access_closed"`
	} `json:"response"`
}


func (rb *RollBot) FindUser(userId int) (string, error) {
	method := "users.get"
	params := make(map[string]string)
	params["user_ids"] = strconv.Itoa(userId)
	response, err := SendMessage(method, params, rb.Config)
	if err != nil {
		return "", err
	}
	u := UserResponse{}
	err = json.Unmarshal(response, &u)
	if err != nil {
		return "", err
	}
	if len(u.Response) < 1 {
		return "", errors.New("user not found")
	}
	name := fmt.Sprintf("%s %s", u.Response[0].FirstName, u.Response[0].LastName)
	return name, nil
}

