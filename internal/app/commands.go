package app

import (
	"strings"
)

type Command uint8

const (
	UnknownCommand Command = iota
	Roll
	CreateCharacter
	Help
)

var CommandHelp = map[Command]string{
	Roll:            "/roll XdY+Z XdY+Z ... XdY+Z - бросок кубиков.\nX - количество, Y - число граней\nZ - дополнительный плюс к результату\nВсе параметры опциональны.",
	CreateCharacter: "/create character - создать две пары аттрибутов для генерации персонажа",
	Help:            "/help - просмотр этой подсказки.",
}

func ParseCommand(data string) (Command, []string) {
	s := strings.Replace(string(data), "\"", "", -1)
	args := strings.Split(s, " ")
	com := strings.ToLower(args[0])
	switch com {
	case "/roll", "/ролл":
		arguments := make([]string, len(args[1:]))
		for i, val := range args[1:] {
			arguments[i] = strings.ToLower(val)
		}
		return Roll, args[1:]
	case "/help", "/h":
		return Help, args[1:]
	case "/create":
		switch strings.ToLower(args[1]) {
		case "character":
			return CreateCharacter, args[2:]
		}
	default:
		return Help, args[1:]
	}
	return UnknownCommand, nil
}
