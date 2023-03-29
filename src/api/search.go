package api

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/tidwall/gjson"

	"github.com/Rosya-edwica/go-hh-parser-vacancies/src/logger"
	"github.com/Rosya-edwica/go-hh-parser-vacancies/src/models"
	"github.com/Rosya-edwica/go-hh-parser-vacancies/src/mysql"
)

const GroupSize = 1

func GetVacanciesByQuery(city models.City, profession models.Profession) {
	url := CreateLink(profession.Name, city.HH_ID)
	checkCaptcha(url)
	json, err := GetJson(url)
	if err != nil {
		logger.Log.Printf("Ошибка при подключении к странице с вакансиями: %s. Error: %s", err, url)
		return
	}
	pagesCount := gjson.Get(json, "pages").Int()
	found := gjson.Get(json, "found").Int()
	if found > 2000 && city.Name == "Russia" {
		logger.Log.Printf("Профессия: %s | Город: %s | Найдено: %d", profession.Name, city.Name, found)
		logger.Log.Printf("Найдено вакансий свыше 2000. Поэтому будет вестись поиск по отдельным городам для этой профессии")
		parseProfessionByCurrentCity(profession)
		return
	} else {
		logger.Log.Printf("Профессия: %s | Город: %s | Найдено: %d", profession.Name, city.Name, found)
		for page := 0; page < int(pagesCount); page++ {
			ParseVacanciesFromPage(fmt.Sprintf("%s&page=%d", url, page), city.EDWICA_ID, profession.Id)
		}
		return
	}
}

func ParseVacanciesFromPage(url string, city_edwica int, id_profession int) {
	json, err := GetJson(url)
	if err != nil {
		logger.Log.Printf("Не удалось подключиться к странице %s.\nТекст ошибки: %s", err, url)
		return
	}

	items := gjson.Get(json, "items").Array()
	var wg sync.WaitGroup
	wg.Add(len(items))
	for _, item := range items {
		go scrapeVacancy(item.Get("url").String(), city_edwica, id_profession, &wg)
	}
	wg.Wait()
	return
}

func parseProfessionByCurrentCity(profession models.Profession) {
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
		GetVacanciesByQuery(city, profession)
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
