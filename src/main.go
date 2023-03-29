package main

import (
	"fmt"
	"time"

	"github.com/Rosya-edwica/go-hh-parser-vacancies/src/api"
	"github.com/Rosya-edwica/go-hh-parser-vacancies/src/logger"
	"github.com/Rosya-edwica/go-hh-parser-vacancies/src/models"
	"github.com/Rosya-edwica/go-hh-parser-vacancies/src/mysql"
)

const GroupSize = 1

func main() {
	logger.Log.Printf("Старт программы")
	start := time.Now().Unix()
	UpdateCurrency()

	Run()

	logger.Log.Println("Время выполнения программы в секундах:", time.Now().Unix()-start)
	fmt.Println("Время выполнения программы в секундах:", time.Now().Unix()-start)

}

func UpdateCurrency() {
	var confirmation string
	fmt.Printf("Обновить текущие значения валюты в БД? [Y/n] ")
	fmt.Scan(&confirmation)
	if confirmation != "Y" {
		logger.Log.Printf("Отмена обновления валюты в БД")
		return
	}

	currency := api.GetCurrencies()
	mysql.UpdateCurrencyRate(currency)
	logger.Log.Print("Валюта в БД была обновлена")
}

func Run() {
	defaultCity := models.City{
		HH_ID:     0,
		EDWICA_ID: 1,
		Name:      "Russia",
	}
	professions := mysql.GetProfessions()
	for _, prof := range professions {
		logger.Log.Printf("Ищем профессию `%s`", prof.Name)
		api.GetVacanciesByQuery(defaultCity, prof)
		mysql.SetParsedStatusToProfession(prof.Id)
		logger.Log.Printf("Профессия %s спарсена", prof.Name)

	}
}
