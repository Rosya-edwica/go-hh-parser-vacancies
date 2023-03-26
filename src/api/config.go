package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/tidwall/gjson"

	"github.com/Rosya-edwica/go-hh-parser-vacancies/src/logger"
	"github.com/Rosya-edwica/go-hh-parser-vacancies/src/models"
	"github.com/Rosya-edwica/go-hh-parser-vacancies/src/telegram"
)

const (
	TOKEN           = "QQAVSIBVU4B0JCR296THKB22JP05A92H329U49TDD9CRIS8DT9BRPPT7M9OLQ6HD"
	dictionariesUrl = "https://api.hh.ru/dictionaries"
	domain          = "https://api.hh.ru/vacancies?"
	per_page        = "100"
	search_field    = "name"
)

var professional_roles = map[int][]string{
	85: {"159", "21", "31", "39", "52", "58", "63", "67", "81", "172", "131"},
	84: {"4", "5", "62", "70"},
}

var industries = map[int][]string{
	86: {"13.662", "388.509"},
}

var headers = map[string]string{
	"User-Agent":    "Mozilla/5.0 (iPad; CPU OS 7_2_1 like Mac OS X; en-US) AppleWebKit/533.14.6 (KHTML, like Gecko) Version/3.0.5 Mobile/8B116 Safari/6533.14.6",
	"Authorization": fmt.Sprintf("Bearer %s", TOKEN),
}
var RelevantCurrencies = []models.Currency{}

func init() {
	RelevantCurrencies = GetCurrencies()
}

func CreateLink(name string, area int, role int) (link string) {
	if role == 86 {
		params := url.Values{
			"search_field": {search_field},
			"per_page":     {per_page},
			"text":         {name},
			"area":         {strconv.Itoa(area)},
			"industry":     industries[role],
		}
		link = domain + params.Encode()
		return
	} else if role == 85 || role == 84{
		params := url.Values{
			"search_field":      {search_field},
			"per_page":          {per_page},
			"text":              {name},
			"area":              {strconv.Itoa(area)},
			"professional_role": professional_roles[role],
		}
		link = domain + params.Encode()
		return
	} else {
		params := url.Values{
			"search_field":      {search_field},
			"per_page":          {per_page},
			"text":              {name},
		}
		link = domain + params.Encode()
		return
	}

}

func GetJson(url string) (json string, err error) {
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	checkErr(err)
	req.Header.Set("User-Agent", headers["User-Agent"])
	req.Header.Set("Authorization", headers["Authorization"])
	response, err := client.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	return string(data), nil
}

func convertSalaryToRUR(salary models.Salary) models.Salary {
	for _, cur := range RelevantCurrencies {
		if cur.Code == salary.Currency {
			salary.To = salary.To / cur.Rate
			salary.From = salary.From / cur.Rate
			salary.Currency = "RUR"
			return salary
		}
	}
	return salary
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func checkCaptcha(url string) {
	json, err := GetJson(url)
	if err != nil {
		logger.Log.Printf("Ошибка при подключении к странице %s.\nТекст ошибки: %s", err, url)
	}
	if checkManyRequestsError(json) {
		captcha_url := gjson.Get(json, "errors.0.captcha_url").String()
		count_updates := telegram.GetUpdates()
		if captcha_url == "" {
			return
		}

		telegram.Mailing(fmt.Sprintf("Пройди капчу по этому адресу: %s&backurl=https://edwica.ru \nПосле отправь мне любое сообщение...\n", captcha_url))
		for {
			fmt.Println("Ждем ответа!!")
			time.Sleep(5 * time.Second)
			new_updates := telegram.GetUpdates()
			if new_updates != count_updates {
				break
			}
		}
		json, err = GetJson(url)
		if err != nil {
			logger.Log.Printf("Ошибка при подключении к странице %s.\nТекст ошибки: %s", err, url)
		}
		if checkManyRequestsError(json) {
			logger.Log.Printf("Не смогли спарсить вакансию: %s\nТекст ошибки: %s", url, json)
			return
		}
	}
}

func checkManyRequestsError(json string) bool {
	err := gjson.Get(json, "errors.0.type").String()
	if err == "" {
		return false
	} else {
		return true
	}
}

func checkErr(err error) {
	if err != nil {
		telegram.Mailing(fmt.Sprintf("Программа остановилась: %s", err))
		logger.Log.Fatal(err)
		panic(err)
	}
}
