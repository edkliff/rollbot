package app

import (
	"errors"
	"fmt"
	"github.com/edkliff/rollbot/internal/generator"
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

func (app *RollBot) RollCommand(a ...string)( string, error) {
	result, err := app.Generator.Roll(2, 6)
	if err != nil {
		return "", err
	}
	resultString := fmt.Sprintf("2d6, результат: %d - %v", generator.Sum(result), result)
	return resultString, nil
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
	return nil, nil, errors.New("что-то произошло непонятное")
}
