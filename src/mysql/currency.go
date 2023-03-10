package mysql

import (
	"fmt"
	"strings"

	"github.com/Rosya-edwica/go-hh-parser-vacancies/src/logger"
	"github.com/Rosya-edwica/go-hh-parser-vacancies/src/models"
)

func UpdateCurrencyRate(currencies []models.Currency) {
	db := connect()
	defer db.Close()

	for _, cur := range currencies {
		query := fmt.Sprintf(`UPDATE %s SET rate=%f WHERE code="%s";`,
			TableCurrencies, cur.Rate, cur.Code)
		tx, _ := db.Begin()
		_, err := db.Exec(query)
		if err == nil {
			tx.Commit()
		} else {
			tx.Commit()
			db.Close()
			logger.Log.Printf("Не удалось обновить валюты - %s", err)
		}
	}
}

func InsertNew(currencies []models.Currency) {
	db := connect()
	defer db.Close()

	valueStrings := []string{}
	valueArgs := []interface{}{}
	valueInsertCount := 1

	for _, cur := range currencies {
		valueStrings = append(valueStrings, buildPatternInsertValues(4))
		valueInsertCount += 4
		valueArgs = append(valueArgs, cur.Code)
		valueArgs = append(valueArgs, cur.Abbr)
		valueArgs = append(valueArgs, cur.Name)
		valueArgs = append(valueArgs, cur.Rate)
	}
	smt := fmt.Sprintf("INSERT INTO %s (code, abbr, name, rate) VALUES", TableCurrencies)
	smt = fmt.Sprintf(smt+"%s", strings.Join(valueStrings, ","))
	tx, err := db.Begin()

	_, err = db.Exec(smt, valueArgs...)
	if err != nil {
		tx.Commit()
		db.Close()
		logger.Log.Printf("Ошибка: Не удалось добавить валюту в базу - %s", err)
		return
	}
	tx.Commit()
	fmt.Println("Успешно добавили валюту в базу!")
}
