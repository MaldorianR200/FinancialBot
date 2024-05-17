package telegram

import (
	"FinancialBot/currency"
	"FinancialBot/lib/e"
	"FinancialBot/storage"
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
)

const (
	GetExchangeRate = "/get_exchange_rate"
	HelpCmd         = "/help"
	StartCmd        = "/start"
	GetIncome       = "/get_income"
	GetExpense      = "/get_expenses"
	AddIncExp       = "/add_inc_exp"
)

func (p *Processor) doCmd(ctx context.Context, text string, chatID int, username string,
	datetime string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, username)

	parts := strings.Fields(text)
	command := parts[0]

	switch command {
	case GetExchangeRate:
		if len(parts) == 1 {
			return p.tg.SendMessage(chatID, "Введите код валюты для получения курса, например, /get_exchange_rate USD")
		} else if len(parts) == 2 {
			curr := parts[1]
			return p.sendExchangeRate(ctx, chatID, curr)
		} else {
			return p.tg.SendMessage(chatID, "Некорректный формат команды. Используйте /get_exchange_rate <код валюты>")
		}
	case HelpCmd:
		return p.sendHelp(ctx, chatID)
	case StartCmd:
		return p.sendHello(ctx, chatID)
	case GetIncome:
		return p.getIncomes(ctx, chatID, username)
	case GetExpense:
		return p.getExpenses(ctx, chatID, username)
	case AddIncExp:
		if len(parts) != 3 {
			return p.tg.SendMessage(chatID, "Некорректный формат команды. Используйте /add_income <доход> <расход>")
		}
		income := parts[1]
		expense := parts[2]
		return p.addIncExp(ctx, chatID, username, income, expense, datetime)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}

}

//New `INSERT INTO pages (user_name, income, expenses, datetime) VALUES (?, ?)`

func (p *Processor) savePage(ctx context.Context, chatID int, income string, expenses string,
	username string, datetime string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()
	//send := NewMessageSender(chatID, tg)	сокращенный вариант отправки сообшения
	page := &storage.Page{
		UserName: username,
		Income:   income,
		Expenses: expenses,
		DateTime: datetime,
	}

	if err := p.storage.Save(ctx, page); err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}
	return nil
}

func (p *Processor) sendExchangeRate(ctx context.Context, chatID int, curr string) error {
	rate, err := currency.GetExchangeRate(curr)
	fmt.Println("sendExchangeRate", rate)
	if err != nil {
		log.Printf("Ошибка получения курса RUB: %v", err)
		return p.tg.SendMessage(chatID, fmt.Sprintf("Ошибка получения курса RUB: %v", err))
	}

	exchangeRate, err := currency.GetExchangeRate(curr)
	if err != nil {
		return p.tg.SendMessage(chatID, fmt.Sprintf("Ошибка получения курса %s: %v", curr, err))
	}

	if exchangeRate == 0 {
		return p.tg.SendMessage(chatID, fmt.Sprintf("Курс для валюты %s не найден", curr))
	}

	rateToRUB := exchangeRate
	message := fmt.Sprintf("Текущий курс %s к RUB: %f", curr, rateToRUB)
	return p.tg.SendMessage(chatID, message)

}

func (p *Processor) sendHelp(ctx context.Context, chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(ctx context.Context, chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func (p *Processor) getIncomes(ctx context.Context, chatID int, username string) error {
	totalIncome, err := p.storage.GetTotalIncome(ctx, username)
	if err != nil {
		return p.tg.SendMessage(chatID, fmt.Sprintf("Ошибка получения доходов: %v", err))
	}

	return p.tg.SendMessage(chatID, fmt.Sprintf("Общий доход: %f", totalIncome))
}

func (p *Processor) getExpenses(ctx context.Context, chatID int, username string) error {
	totalExpense, err := p.storage.GetTotalExpense(ctx, username)
	if err != nil {
		return p.tg.SendMessage(chatID, fmt.Sprintf("Ошибка получения расходов: %v", err))
	}

	return p.tg.SendMessage(chatID, fmt.Sprintf("Общие расходы: %f", totalExpense))
}

func (p *Processor) addIncExp(ctx context.Context, chatID int, username string, income string, expense string, datetime string) error {
	page := &storage.Page{
		UserName: username,
		Income:   income,
		Expenses: expense,
		DateTime: datetime,
	}

	if err := p.storage.Save(ctx, page); err != nil {
		return p.tg.SendMessage(chatID, fmt.Sprintf("Ошибка сохранения данных: %v", err))
	}

	return p.tg.SendMessage(chatID, "Данные успешно сохранены")
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
