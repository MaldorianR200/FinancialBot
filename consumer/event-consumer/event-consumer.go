package event_consumer

import (
	"FinancialBot/events"
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c Consumer) Start(ctx context.Context) error {
	for {
		// Максимальное кол-во попыток
		maxRetries := 3
		retries := 0
		for retries <= maxRetries {
			gotEvents, err := c.fetcher.Fetch(c.batchSize)
			if err != nil {
				log.Printf("[ERR] consumer: %s ", err.Error())
				retries++
				continue
			}

			if len(gotEvents) == 0 {
				// Ждём 1 секунду
				time.Sleep(1 * time.Second)

				continue
			}

			// Создаём WaitGroup для параллельной обработки событий
			var wg sync.WaitGroup
			wg.Add(len(gotEvents))

			// Канал для передачи ошибок из горутины в основной поток
			errCh := make(chan error, len(gotEvents))

			// Параллельная обработка событий
			for _, event := range gotEvents {
				go func(event events.Event) {
					defer wg.Done()
					if err := c.processor.Process(ctx, event); err != nil {
						errCh <- err
					}
				}(event)
			}

			// Ожидание завершения всех горутин
			wg.Wait()
			close(errCh)

			// Проверка наличия ошибок в горутинах
			for err := range errCh {
				log.Printf("can't handle event: %s ", err.Error())
			}

			//if err := c.handleEvents(ctx, gotEvents); err != nil {
			//	log.Printf(err.Error())
			//
			//	continue
			//}

			// События успешно обработаны, выходим из цикла попыток
			break
		}
		// Проверяем, достигли ли максимального количества попыток
		if retries == maxRetries {
			return errors.New("maximum retries exceeded")
		}
	}
}

/*
	1. Потеря событий: ретраи, возвращение в хранилище, фоллбэк, подтверждение для Featcher
	2. Обработка всей пачки: останавливаться после первой ошибки, счётчик ошибок
	3. Параллельная обработка
*/
// sync.WaitGroup{} <- 3
func (c Consumer) handleEvents(ctx context.Context, events []events.Event) error {
	fmt.Println("Обрабтываем события")
	for _, event := range events {
		log.Printf("got new event: %s", event.Text)

		if err := c.processor.Process(ctx, event); err != nil {
			log.Printf("can't handle event: %s", err.Error())

			continue
		}
	}

	return nil
}
