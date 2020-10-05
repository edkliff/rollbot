package app

import (
	"fmt"
	"github.com/edkliff/rollbot/internal/generator"
	"regexp"
	"strconv"
	"strings"
)

type RollResult struct {
	results []rollResultSet
	comment string
	finalSum int64
}

type rollResultSet struct {
	count int64
	dice int64
	adder int64
	results []int64
	sum int64
}

func (app *RollBot) RollCommand(vk VKReq) (Resulter, error) {
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
	rr := RollResult{
		results: make([]rollResultSet, 0),
		comment: reason,
		finalSum: 0,
	}
	fmt.Println(args)
	if len(args) <2 {
		args = append(args, "1d6+0")
	}
	for _, arg := range args {
		count, dice, adder, isCommand := ParseRoll(arg)
		if !isCommand {
			continue
		}
		res, err := app.Generator.Roll(count, dice)
		if err != nil {
			return nil, err
		}
		r := rollResultSet{
			count:  count,
			dice:   dice,
			adder:  adder,
			results: res,
			sum: generator.Sum(res)+adder,
		}
		rr.results = append(rr.results, r)
	}
	for _, v := range rr.results {
		rr.finalSum += v.sum
	}
	fmt.Println(&rr)
	return &rr, nil
}

func ParseRoll(s string) (int64, int64, int64, bool) {
	ok, err := regexp.Match("^\\d{1,3}d\\d{1,3}[+]\\d{1,3}$", []byte(s))
	if err != nil {
		return 0,0,0, false
	}
	if ok {
		count, dice, adder := SplitRoll(s)
		return count, dice, adder, true
	}

	ok, err = regexp.Match("^d\\d{1,3}[+]\\d{1,3}$", []byte(s))
	if err != nil {
		return 0,0,0, false
	}
	if ok {
		s = "1"+s
		count, dice, adder := SplitRoll(s)
		return count, dice, adder, true
	}

	ok, err = regexp.Match("^\\d{1,3}d\\d{1,3}$", []byte(s))
	if err != nil {
		return 0,0,0, false
	}
	if ok {
		s += "+0"
		count, dice, adder := SplitRoll(s)
		return count, dice, adder, true
	}
	ok, err = regexp.Match("^d\\d{1,3}$", []byte(s))
	if err != nil {
		return 0,0,0, false
	}
	if ok {
		s = "1" + s + "+0"
		count, dice, adder := SplitRoll(s)
		return count, dice, adder, true
	}
	return 0,0,0, false
}

func SplitRoll(s string) (int64, int64, int64)  {
	fst := strings.Split(s, "+")
	sst := strings.Split(fst[0], "d")
	count, err := strconv.Atoi(sst[0])
	if err!= nil {
		return 0,0,0
	}
	dice, err := strconv.Atoi(sst[1])
	if err!= nil {
		return 0,0,0
	}
	adder, err := strconv.Atoi(fst[1])
	if err!= nil {
		return 0,0,0
	}
	return int64(count), int64(dice), int64(adder)
}

func (r *RollResult) VKString() string {
	if len(r.results) == 0 {
		return "Не удалось распарсить ни один из аргументов. Введите в формате XXXdYYY+ZZZ"
	}
	s := fmt.Sprintf("Общая сумма: %d\nПодробно:\n", r.finalSum)
	if r.comment != "" {
		s = fmt.Sprintf("%s\n%s",r.comment, s)
	}
	for _, res := range r.results{
		s += fmt.Sprintf("%dd%d+%d: %d - %v+%d\n",
			res.count, res.dice, res.adder, res.sum, res.results, res.adder)
	}
	return s
}

func (r *RollResult) Comment() string  {
	return r.comment
}

func (r *RollResult) HTML() string {
	if len(r.results) == 0 {
		return "Не удалось распарсить ни один из аргументов. Введите в формате XXXdYYY+ZZZ"
	}
	s := fmt.Sprintf("<div>Сумма: %d</div>", r.finalSum)
	for _, res := range r.results{
		s += fmt.Sprintf("<div>%dd%d+%d: %d - %v+%d\n</div>",
			res.count, res.dice, res.adder, res.sum, res.results, res.adder)
	}
	return ""
}