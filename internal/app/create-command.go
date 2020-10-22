package app

import (
	"fmt"
	"github.com/edkliff/rollbot/internal/generator"
	"html"
	"strings"
)

type CreateResult struct {
	set1 []characteristic
	set2 []characteristic

	sumSet1 int64
	sumSet2 int64

	reason string
}

type characteristic struct {
	sum  int64
	roll []int64
}

func (app *RollBot) CreateCommand(vk VKReq) (Resulter, error) {
	tmpString := vk.Object.Message.Text
	reason, err := GetReason(tmpString)
	if err != nil {
		return nil, err
	}
	if reason != "" {
		tmpString = strings.ReplaceAll(tmpString, reason, "")
	}
	reason = strings.ReplaceAll(strings.ReplaceAll(reason, "(", ""), ")", "")
	set1, err := app.CreateCharacteristics()
	if err != nil {
		return nil, err
	}
	set2, err := app.CreateCharacteristics()
	if err != nil {
		return nil, err
	}
	c := CreateResult{
		set1:    set1,
		set2:    set2,
		sumSet1: 0,
		sumSet2: 0,
		reason:  reason,
	}
	for _, v := range set1 {
		c.sumSet1 += v.sum
	}
	for _, v := range set2 {
		c.sumSet2 += v.sum
	}
	return &c, nil
}

func (app *RollBot) CreateCharacteristics() ([]characteristic, error) {
	resultSet := make([]characteristic, 0, 6)
	for i := 0; i < 6; i++ {
		ok := true
		c := characteristic{}
		for ok {
			res, err := app.Generator.Roll(4, 6)
			if err != nil {
				return nil, err
			}
			res = generator.Sort(res)
			c.roll = res
			c.sum = generator.Sum(res[:3])
			if c.sum >= 8 {
				ok = false
			}
		}
		resultSet = append(resultSet, c)
	}
	return resultSet, nil
}

func (h *CreateResult) VKString() string {
	s := fmt.Sprintf("Набор 1:\nВсего очков: %d\n", h.sumSet1)
	for _, v := range h.set1 {
		s += fmt.Sprintf("%d : %v\n", v.sum, v.roll)
	}
	s += fmt.Sprintf("\nНабор 2:\nВсего очков: %d\n", h.sumSet2)
	for _, v := range h.set2 {
		s += fmt.Sprintf("%d : %v\n", v.sum, v.roll)
	}
	if h.reason != "" {
		s = fmt.Sprintf("%s\n%s", h.reason, s)
	}
	return s
}

func (h *CreateResult) Comment() string {
	return html.EscapeString(h.reason)
}

func (h *CreateResult) HTML() string {
	s := fmt.Sprintf("<div>Набор 1:</div><div>Всего очков: %d</div>", h.sumSet1)
	for _, v := range h.set1 {
		s += fmt.Sprintf("<div>%d : %v</div>", v.sum, v.roll)
	}
	s += fmt.Sprintf("<br><div>Набор 2:</div><div>Всего очков: %d</div>", h.sumSet2)
	for _, v := range h.set2 {
		s += fmt.Sprintf("<div>%d : %v</div>", v.sum, v.roll)
	}
	return s
}
