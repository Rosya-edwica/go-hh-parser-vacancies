package api

import (
	"github.com/tidwall/gjson"

	"github.com/Rosya-edwica/go-hh-parser-vacancies/src/logger"
	"github.com/Rosya-edwica/go-hh-parser-vacancies/src/models"
)

func GetCurrencies() (currencies []models.Currency) {
	json, err := GetJson(dictionariesUrl)
	if err != nil {
		logger.Log.Printf("Не удалось обновить валюту. Текст сообщения: %s", err)
	}
	for _, item := range gjson.Get(json, "currency").Array() {
		currencies = append(currencies, models.Currency{
			Code: item.Get("code").String(),
			Abbr: item.Get("abbr").String(),
			Name: item.Get("name").String(),
			Rate: item.Get("rate").Float(),
		})
	}
	return
}
