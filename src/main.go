package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Rosya-edwica/go-hh-parser-vacancies/src/api"
	"github.com/Rosya-edwica/go-hh-parser-vacancies/src/logger"
	"github.com/Rosya-edwica/go-hh-parser-vacancies/src/models"
	"github.com/Rosya-edwica/go-hh-parser-vacancies/src/mysql"
	"github.com/joho/godotenv"
)

const GroupSize = 1
var ProfessionsFromFile bool

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Нет файла с переменными окружения .env!")
	}
	val := os.Getenv("PROFESSIONS_FROM_FILE")
	if val == "true" {
		ProfessionsFromFile = true
	} else {
		ProfessionsFromFile = false
	}
	}	

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
	professions := mysql.GetProfessionsFromFile(ProfessionsFromFile)
	for _, profession := range professions {
		logger.Log.Printf("Ищем профессию `%s`", profession.Name)
		profession.OtherNames = append(profession.OtherNames, profession.Name)
		unique_professions := api.Unique_list(profession.OtherNames)
		for _, prof := range unique_professions {
			if len(prof) <= 3 {
				continue
			}
			api.GetVacanciesByQuery(defaultCity, profession)

		}
		// mysql.SetParsedStatusToProfession(profession.Id)
		logger.Log.Printf("Профессия %s спарсена", profession.Name)

	}
}
