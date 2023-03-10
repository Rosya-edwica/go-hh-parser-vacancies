package telegram

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Rosya-edwica/go-hh-parser-vacancies/src/logger"
)

const token = "6105028983:AAG5qYpp0KKkhBiyOHri6MmB5e_UVzfC9pU"

var chats = []string{
	"544490770",
}

func getJson(url string) (json string, err error) {
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	checkErr(err)
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

func checkErr(err error) {
	if err != nil {
		Mailing(fmt.Sprintf("Программа остановилась: %s", err))
		logger.Log.Fatal(err)
		panic(err)
	}
}

func getUrl() string {
	return fmt.Sprintf("https://api.telegram.org/bot%s", token)
}
