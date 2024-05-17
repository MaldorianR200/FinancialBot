package currency

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Тип для хранения данных о курсе валюты
type ExchangeRate struct {
	Base string `json:"base"`
	Date string `json:"date"`
	//Rates struct {
	//	USD float64 `json:"USD"`
	//	RUB float64 `json:"RUB"`
	//	EUR float64 `json:"EUR"`
	//	CNY float64 `json:"CNY"`
	//} `json:"rates"`
	Rates map[string]float64 `json:"rates"`
}

const (
	apiKey          = "ddc31d0d81c9d1f8d1a1bddf"
	exchangeRateAPI = "https://api.exchangerate-api.com/v4/latest/"
)

// Метод для получения текущего курса валюты
func GetExchangeRate(currency string) (float64, error) {

	// Создаем HTTP-клиент
	client := http.Client{}

	// Формируем URL для запроса
	exchangeRateURL := exchangeRateAPI + currency

	// Создаем запрос
	req, err := http.NewRequest("GET", exchangeRateURL, nil)
	if err != nil {
		return 0, err
	}

	// Добавляем API ключ в заголовок запроса
	req.Header.Add("X-API-KEY", apiKey)

	// Отправляем запрос и получаем ответ
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var exchangeRate ExchangeRate
	if err := json.NewDecoder(resp.Body).Decode(&exchangeRate); err != nil {
		return 0, err
	}

	rate, ok := exchangeRate.Rates["RUB"]
	if !ok {
		return 0, errors.New("курс RUB не найден")
	}
	fmt.Println("RATE - ", rate)
	return rate, nil
}
