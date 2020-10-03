package app

type Command uint8

const (
	UnknownCommand Command = iota
	Roll
	CreateCharacter
	Help
)

func CommandHelp(a ...[]string)( string, error ){
	return "/roll XdY+Z XdY+Z ... XdY+Z REASON - бросок кубиков.\n" +
		"X - количество, Y - число граней\n" +
		"Z - дополнительный плюс к результату, REASON - описание броска\n" +
		"Все параметры опциональны.\n"+
 		"/create character - создать две пары аттрибутов для генерации персонажа\n"+
    	"/help - просмотр этой подсказки.", nil
}

func RollCommand(a ...[]string)( string, error) {
	return "", nil
}