package app

import (
	"errors"
	"regexp"
	"strings"
)

type Resulter interface {
	HTML() string
	VKString() string
	Comment() string
}

type ErrorResult struct {
	err error
}

func NewErrorResult(err error) *ErrorResult {
	return &ErrorResult{err}
}

func (h *ErrorResult) VKString() string {
	return h.err.Error()
}

func (h *ErrorResult) Comment() string {
	return h.err.Error()
}

func (h *ErrorResult) HTML() string {
	return "<div>" + h.err.Error() + "</div>"
}

func (app *RollBot) ParseCommand(vkr *VKReq) (func(VKReq) (Resulter, error), error) {
	vkr.Object.Message.Text = strings.TrimSpace(vkr.Object.Message.Text)
	argsList := strings.Split(vkr.Object.Message.Text, " ")
	if len(argsList) > 0 {
		switch strings.ToLower(argsList[0]) {
		case "/roll", "/r":
			return app.RollCommand, nil
		case "/help", "/h":
			return app.HelpCommand, nil
		case "/create", "/c":
			return app.CreateCommand, nil
		case "/%":
			return app.PercentCommand, nil
		}
	}
	return nil, errors.New("unknown command")
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
	reason := string(reasonb)
	return reason, nil
}
