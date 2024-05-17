package main

import (
	tgClient "FinancialBot/clients/telegram"
	"FinancialBot/consumer/event-consumer"
	"FinancialBot/events/telegram"
	"FinancialBot/storage/sqlite"
	"context"
	"flag"
	"log"
)

const (
	tgBotHost         = "api.telegram.org"
	sqliteStoragePath = "data/sqlite/storage.db"
	batchSize         = 100
)

/*
	Финансовый бот: Может предоставлять информацию о курсах валют,
	биржевых котировках, помогать в отслеживании расходов и доходов.
*/

// botToken := "6868760353:AAGPJzkReqYKX0XdxbePeSmxvRKtrVYjPAI"
func main() {
	// s := files.New(storagePath)
	ctx := context.Background()
	s, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Fatal("can't connect to sqlite storage: ", err)
	}
	// context.Background() - данный контекст ни в чём нас не ограничивает
	if err := s.Init(context.TODO()); err != nil {
		log.Fatal("can't init sqlite storage: ", err)
	}
	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		s,
	)
	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(ctx); err != nil {
		log.Fatal(err)
	}
}

// приставка must говорит о том,  что с данной функцие стоит работать осторожно
func mustToken() string {
	token := flag.String(
		"tg-bot-token", // название флага - name
		"6868760353:AAGPJzkReqYKX0XdxbePeSmxvRKtrVYjPAI", // значение, которое будет здесь храниться - value
		"token for access to telegram bot",               // описание флага
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
