package telegram

const msgHelp = `Я могу предоставить вам курс валюты на данный момент. 
Также я могу предложить вам отслеживать свои доходы и расходы.
Чтобы узнать курск валюты введите команду /GetExchangeRate. А если вы хотите сохранить
ваши доходы и расходы введите команду /money
`

const msgHello = "Привет! \n\n" + msgHelp

const (
	msgUnknownCommand = "Неизвестная команда🤔" // Unknown command

	msgSaved = "Saved!👍"
)

//const msgHelp = `I can save and keep you pages. Also I can offer you them to read.
//In order to save the page, just send me al link ti it.
//
//In order to get a random page from your list, send me command /rnd.
//Caution! After that, this page will be removed from your list!`
