package app

import (
	"errors"
	"fmt"
	"github.com/edkliff/rollbot/internal/generator"
	"regexp"
	"strconv"
	"strings"
)

type Command uint8

const (
	UnknownCommand Command = iota
	Roll
	CreateCharacter
	Help
)

func (app *RollBot) HelpCommand(a ...string)( string, error ){
	return "/roll XdY+Z XdY+Z ... XdY+Z REASON - бросок кубиков.\n" +
		"X - количество, Y - число граней\n" +
		"Z - дополнительный плюс к результату, REASON - описание броска\n" +
		"Все параметры опциональны.\n"+
 		"/create character - создать две пары аттрибутов для генерации персонажа\n"+
    	"/help - просмотр этой подсказки.", nil
}

func (app *RollBot) RollCommand(args ...string)( string, error) {
	if len(args) == 0 {
		result, err := app.Generator.Roll(2, 6)
		if err != nil {
			return "", err
		}
		resultString := fmt.Sprintf("2d6: %d - %v", generator.Sum(result), result)
		return resultString, nil
	}
	args, reason, err := GetReason(strings.Join(args, " "))
	if err != nil {
		return "", err
	}
	resultString := ""
	finalSum := int64(0)
	for _, arg := range args {
		if isRoll(arg) {
			count, dice, adder, err := ParseRoll(arg)
			if err != nil {
				return "", err
			}
			r, err := app.Generator.Roll(count, dice)
			if err != nil {
				return "", err
			}
			sum := generator.Sum(r) + adder
			finalSum += sum
			resultString+=fmt.Sprintf("%dd%d+%d: %d - %v", count, dice, adder, sum, r)
		}
	}
	if len(reason) > 0 {
		resultString = reason + "\n"
	}

	resultString = fmt.Sprintf("Сумма: %d\n", finalSum) + resultString
	return resultString, nil
}

func GetReason(s string) ([]string, string, error) {
	seq := "\\(.*\\)"
	rx, err := regexp.Compile(seq)
	if err != nil {
		return make([]string, 0), "", err
	}
	reasonb := rx.Find([]byte(s))
	if reasonb == nil {
		args := strings.Split(s, " ")
		return args, "", nil
	}
	reason := string(reasonb)
	withoutReason := strings.ReplaceAll(s, reason, "")
	args := strings.Split(withoutReason, " ")
	return args, reason, nil
}

func isRoll(s string) bool {
	 ok, err := regexp.Match("[0-9]*d?[0-9]*\\+?[0-9]*", []byte(s))
	 if err != nil {
		return false
	}
	return ok
}

func ParseRoll(s string) (int64, int64, int64, error)  {
	seq := "[0-9]*d?[0-9]*\\+?[0-9]*"
	rx, err := regexp.Compile(seq)
	if err != nil {
		return 0,0,0, err
	}
	sb := rx.Find([]byte(s))
	if sb == nil {
		return 0,0,0, errors.New("not finded")
	}
	s = string(sb)
	s = strings.ReplaceAll(s, "+", " ")
	s = strings.ReplaceAll(s, "d", " ")
	args := strings.Split(s, " ")
	fmt.Println(args)
	var count, dice, adder int64
	if len(args) >=1 {
		counts := args[0]
		counti, err := strconv.Atoi(counts)
		if err != nil && counts != "" {
			return 0,0,0, errors.New("Малой - мудак")
		}
		if counts == "" {
			counti = 1
		}
		count = int64(counti)
	}
	if len(args)  >=2{
		dices := args[1]
		decei, err := strconv.Atoi(dices)
		if err != nil {
			return 0,0,0, errors.New("Малой - мудак")
		}
		if dices == "" {
			decei = 6
		}
		dice = int64(decei)
	}
	if len(args) >=3 {
		adders := args[2]
		adderi, err := strconv.Atoi(adders)
		if err != nil {
			return 0,0,0, errors.New("Малой - мудак")
		}
		adder = int64(adderi)
	}


	return count, dice, adder, nil
}

func (app *RollBot) ParseCommand(vkr *VKReq)(func(...string)(string, error), []string, error) {
	vkr.Object.Message.Text = strings.TrimSpace(strings.ToLower(vkr.Object.Message.Text))
	argsList := strings.Split(vkr.Object.Message.Text, " ")
	if len(argsList) > 0 {
		args := make([]string, 0)
		if len(argsList) > 1 {
			args = argsList[1:]
		}
		switch argsList[0] {
		case "/roll":
			return app.RollCommand, args, nil
		case "/help":
			return app.HelpCommand, args, nil
		}
	}
	return app.HelpCommand, make([]string, 0), nil
}
