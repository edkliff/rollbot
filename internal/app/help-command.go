package app

import "strings"

type HelpResult struct {
	helps []string
	reason string
}

func (app *RollBot) HelpCommand(vk *VKReq) (Resulter, error) {
	h := HelpResult{}
	h.helps = []string{"/roll XdY+Z XdY+Z ... XdY+Z REASON - бросок кубиков.\n",
		"X - количество, Y - число граней\n",
		"Z - дополнительный плюс к результату, REASON - описание броска\n",
		"Все параметры опциональны.\n",
		"/help - просмотр этой подсказки.",
	}
	comment, err := GetReason(vk.Object.Message.Text)
	if err !=nil {
		return nil, err
	}
	h.reason = comment
	return  &h, nil
}

func (h *HelpResult) String() string {
	s := strings.Join(h.helps, "")
	return s
}

func (h *HelpResult) Comment() string  {
	return h.reason
}

func (h *HelpResult) HTML() string {
	s := strings.Join(h.helps, "</div><div>")
	return "<div>" + s + "</div>"
}