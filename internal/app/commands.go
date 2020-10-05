package app

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Resulter interface {
	HTML() string
	fmt.Stringer
	Comment() string
}

type ErrorResult struct {
	err error
}

func NewErrorResult(err error) *ErrorResult {
	return &ErrorResult{err}
}

func (h *ErrorResult) String() string {
	return h.err.Error()
}

func (h *ErrorResult) Comment() string  {
	return h.err.Error()
}

func (h *ErrorResult) HTML() string {
	return "<div>" + h.err.Error() + "</div>"
}

func (app *RollBot) ParseCommand(vkr *VKReq) (func(*VKReq) (Resulter, error), error) {
	vkr.Object.Message.Text = strings.TrimSpace(strings.ToLower(vkr.Object.Message.Text))
	argsList := strings.Split(vkr.Object.Message.Text, " ")
	if len(argsList) > 0 {
		switch argsList[0] {
		case "/roll":
			return app.RollCommand,  nil
		case "/help":
			return app.HelpCommand, nil
		}
	}
	return app.HelpCommand, nil
}


func GetReason(s string) (string, error) {
	seq := "\\(.*\\)"
	rx, err := regexp.Compile(seq)
	if err != nil {
		return "", err
	}
	reasonb := rx.Find([]byte(s))
	if reasonb == nil {
		return "", nil
	}
	reason := string(reasonb[1:len(reasonb)-1])
	return reason, nil
}

func ParseRoll(s string) (int64, int64, int64, error) {
	seq := "[0-9]*d?[0-9]*\\+?[0-9]*"
	rx, err := regexp.Compile(seq)
	if err != nil {
		return 0, 0, 0, err
	}
	sb := rx.Find([]byte(s))
	if sb == nil {
		return 0, 0, 0, errors.New("not finded")
	}
	s = string(sb)
	s = strings.ReplaceAll(s, "+", " ")
	s = strings.ReplaceAll(s, "d", " ")
	args := strings.Split(s, " ")
	count, dice, adder := int64(1), int64(6), int64(0)
	if len(args) >= 1 {
		counts := args[0]
		counti, err := strconv.Atoi(counts)
		if err != nil && counts != "" {
			return 0, 0, 0, errors.New("Малой - мудак")
		}
		if counts == "" {
			counti = 1
		}
		count = int64(counti)
	}
	if len(args) >= 2 {
		dices := args[1]
		decei, err := strconv.Atoi(dices)
		if err != nil {
			return 0, 0, 0, errors.New("Малой - мудак")
		}
		if dices == "" {
			decei = 6
		}
		dice = int64(decei)
	}
	if len(args) >= 3 {
		adders := args[2]
		adderi, err := strconv.Atoi(adders)
		if err != nil {
			return 0, 0, 0, errors.New("Малой - мудак")
		}
		adder = int64(adderi)
	}


	return count, dice, adder, nil
}
