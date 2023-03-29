package api

import (
	"fmt"
	"strings"
	"sync"

	"github.com/tidwall/gjson"

	"github.com/Rosya-edwica/go-hh-parser-vacancies/src/logger"
	"github.com/Rosya-edwica/go-hh-parser-vacancies/src/models"
	"github.com/Rosya-edwica/go-hh-parser-vacancies/src/mysql"
)

func scrapeVacancy(url string, city_edwica int, id_profession int, wg *sync.WaitGroup) {
	var vacancy models.Vacancy
	// checkCaptcha(url)
	json, err := GetJson(url)
	if err != nil {
		logger.Log.Printf("Ошибка при подключении к странице %s.\nТекст ошибки: %s", err, url)
		return
	}

	salary := getSalary(json)
	vacancy.CityId = city_edwica
	vacancy.SalaryFrom = salary.From
	vacancy.ProfessionId = id_profession
	vacancy.SalaryTo = salary.To
	vacancy.Skills = getSkills(json)
	vacancy.Specializations = getSpecializations(json)
	vacancy.ProfAreas = getProfAreas(json)
	vacancy.Id = int(gjson.Get(json, "id").Int())
	vacancy.Title = gjson.Get(json, "name").String()
	vacancy.Url = gjson.Get(json, "alternate_url").String()
	vacancy.Experience = gjson.Get(json, "experience.name").String()
	vacancy.DateUpdate = gjson.Get(json, "created_at").String()
	mysql.SaveOneVacancy(vacancy)
	wg.Done()
}

func getSalary(vacancyJson string) (salary models.Salary) {
	salary.Currency = gjson.Get(vacancyJson, "salary.currency").String()
	salary.From = gjson.Get(vacancyJson, "salary.from").Float()
	salary.To = gjson.Get(vacancyJson, "salary.to").Float()

	switch salary.Currency {
	case "RUR":
		return salary
	case "":
		return models.Salary{}
	default:
		return convertSalaryToRUR(salary)
	}

}

func getSpecializations(vacancyJson string) string {
	var specializations []string
	for _, item := range gjson.Get(vacancyJson, "specializations").Array() {
		specializations = append(specializations, item.Get("name").String())
	}
	return strings.Join(removeDuplicateStr(specializations), "|")
}

func getProfAreas(vacancyJson string) string {
	var profAreas []string
	for _, item := range gjson.Get(vacancyJson, "specializations").Array() {
		profAreas = append(profAreas, item.Get("profarea_name").String())
	}
	return strings.Join(removeDuplicateStr((profAreas)), "|")
}

func getSkills(vacancyJson string) string {
	var skills []string
	for _, item := range gjson.Get(vacancyJson, "key_skills").Array() {
		skills = append(skills, item.Get("name").String())
	}
	languages := getLanguages(vacancyJson)
	skills = append(skills, languages...)
	return strings.Join(skills, "|")
}

func getLanguages(vacancyJson string) (languages []string) {
	for _, item := range gjson.Get(vacancyJson, "languages").Array() {
		lang := item.Get("name").String()
		level := item.Get("level.name").String()
		languages = append(languages, fmt.Sprintf("%s (%s)", lang, level))
	}
	return
}
