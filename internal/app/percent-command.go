package app

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type PercentResult struct {
	results []PResult
	reason  string
}

type PResult struct {
	percent string
	result  bool
}

func (app *RollBot) PercentCommand(vk VKReq) (Resulter, error) {
	h := PercentResult{
		results: make([]PResult, 0),
		reason:  "",
	}
	tmpString := vk.Object.Message.Text
	reason, err := GetReason(tmpString)
	if err != nil {
		return nil, err
	}
	if reason != "" {
		tmpString = strings.ReplaceAll(tmpString, reason, "")
	}
	reason = strings.ReplaceAll(strings.ReplaceAll(reason, "(", ""), ")", "")
	args := strings.Split(tmpString, " ")
	for _, val := range args {
		ok, err := regexp.Match("^\\d{1,3}%$", []byte(val))
		if err != nil {
			return nil, err
		}
		if ok {
			pr := PResult{
				percent: val,
				result:  false,
			}
			val = strings.ReplaceAll(val, "%", "")
			p, err := strconv.Atoi(val)
			if err != nil {
				return nil, err
			}
			p = 100 - p
			r, err := app.Generator.Roll(1, 100)
			if err != nil {
				return nil, err
			}
			if r[0] > int64(p) {
				pr.result = true
			}
			h.results = append(h.results, pr)
		}

	}
	return &h, nil
}

func (h *PercentResult) VKString() string {
	s := ""
	for _, v := range h.results {
		if v.result {
			s += fmt.Sprintf("%s : Успех\n", v.percent)
			continue
		}
		s += fmt.Sprintf("%s : Провал\n", v.percent)
	}
	return s
}

func (h *PercentResult) Comment() string {
	return h.reason
}

func (h *PercentResult) HTML() string {
	s := ""
	for _, v := range h.results {
		if v.result {
			s += fmt.Sprintf("<div>%s : Успех</div>", v.percent)
			continue
		}
		s += fmt.Sprintf("<div>%s : Провал</div>", v.percent)
	}
	return s
}
