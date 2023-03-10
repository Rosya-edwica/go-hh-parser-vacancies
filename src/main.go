package main

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/Rosya-edwica/go-hh-parser-vacancies/src/api"
	"github.com/Rosya-edwica/go-hh-parser-vacancies/src/logger"
	"github.com/Rosya-edwica/go-hh-parser-vacancies/src/models"
	"github.com/Rosya-edwica/go-hh-parser-vacancies/src/mysql"
)

const GroupSize = 10

func main() {
	logger.Log.Printf("Старт программы")
	start := time.Now().Unix()
	UpdateCurrency()
	professions := mysql.GetProfessions()
	for _, prof := range professions {
		parseProfession(prof)
	}
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

func parseProfession(profession models.Profession) {
	logger.Log.Printf("Ищем профессию `%s`", profession.Name)
	groups := groupCities()
	for _, group := range groups {
		var wg sync.WaitGroup
		wg.Add(len(group))
		for _, city := range group {
			go parseProfessionInCity(city, profession, &wg)
		}
		wg.Wait()
	}
	mysql.SetParsedStatusToProfession(profession.Id)
	logger.Log.Printf("Профессия %s спарсена", profession.Name)

}

func groupCities() (groups [][]models.City) {
	cities := mysql.GetCities()
	citiesCount := len(cities)
	var limit int
	for i := 0; i < citiesCount; i += GroupSize {
		limit += GroupSize
		if limit > citiesCount {
			limit = citiesCount
		}
		group := cities[i:limit]
		groups = append(groups, group)
	}
	logger.Log.Printf("Ведем поиск профессии в  %d городах одновременно", GroupSize)
	return
}

func parseProfessionInCity(city models.City, profession models.Profession, wg *sync.WaitGroup) {
	defer wg.Done()
	profession.OtherNames = append(profession.OtherNames, profession.Name)
	unique_professions := unique_list(profession.OtherNames)
	for _, prof := range unique_professions {
		if len(prof) <= 3 {
			continue
		}
		api.GetVacanciesByQuery(city, prof, profession.Id, profession.ProfRoleId)
	}
}

func unique_list(list []string) []string {
	var unique []string
	var re = regexp.MustCompile(`(?m) +`)

	for _, item := range list {
		no_symbol := strings.ReplaceAll(item, "-", " ")
		trim := re.ReplaceAllString(no_symbol, " ")
		low := strings.ToLower(trim)
		has_match := false
		for _, item2 := range unique {
			if item2 == low {
				has_match = true
				break
			}
		}
		if !has_match {
			unique = append(unique, low)
		}
	}
	return unique
}
