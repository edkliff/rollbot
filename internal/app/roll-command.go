package app

import (
	"fmt"
	"regexp"
	"strings"
)

type RollResult struct {

}

func (app *RollBot) RollCommand(vk *VKReq) (Resulter, error) {
	tmpString := strings.TrimSpace(vk.Object.Message.Text)
	args := strings.Split(tmpString, " ")
	fmt.Println(args)
	return &RollResult{}, nil
}

func isRoll(s string) bool {
	ok, err := regexp.Match("[0-9]*d?[0-9]*\\+?[0-9]*", []byte(s))
	if err != nil {
		return false
	}
	return ok
}


func (h *RollResult) String() string {
	return ""
}

func (h *RollResult) Comment() string  {
	return ""
}

func (h *RollResult) HTML() string {
	s := ""
	return "<div>" + s + "</div>"
}