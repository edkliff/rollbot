package app

import (
	"fmt"
	"regexp"
	"strings"
)

type RollResult struct {
	results [][]int64

}

func (app *RollBot) RollCommand(vk *VKReq) (Resulter, error) {
	tmpString := strings.TrimSpace(vk.Object.Message.Text)
	reason, err := GetReason(tmpString)
	if err != nil {
		return nil, err
	}
	tmpString = strings.ReplaceAll(tmpString, reason, "")
	reason = strings.ReplaceAll(strings.ReplaceAll(reason, "(", ""), ")", "")
	args := strings.Split(tmpString, " ")
	fmt.Println(args)
	if len(args) <2 {
		// r, err := app.Generator.Roll(1, 6)
		// if err != nil {
		// 	return nil, err
		// }

	}
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

	return ""
}