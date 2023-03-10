package mysql

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/Rosya-edwica/go-hh-parser-vacancies/src/logger"
	"github.com/Rosya-edwica/go-hh-parser-vacancies/src/telegram"

	_ "github.com/go-sql-driver/mysql"
)

var (
	TableCity        = "h_city"
	TableProfessions = "h_position"
	TableCurrencies  = "h_currency"
	TableVacancy     = "h_vacancy"
	host             string
	port             string
	user             string
	password         string
	dbname           string
)

func init() {
	host = os.Getenv("MYSQL_HOST")
	port = os.Getenv("MYSQL_PORT")
	user = os.Getenv("MYSQL_USER")
	password = os.Getenv("MYSQL_PASSWORD")
	dbname = os.Getenv("MYSQL_DATABASE")
	if host == "" {
		logger.Log.Fatal("Проверь переменные окружения (MYSQL_HOST, MYSQL_PORT, MYSQL_PASSWORD, MYSQL_USER, MYSQL_DATABASE)")
	}
}

func connect() (db *sql.DB) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbname))
	checkErr(err)
	return db
}

// returned string like -> (?, ?, ?, ?, ...., valuesCount)
func buildPatternInsertValues(valuesCount int) (pattern string) {
	var items []string
	for i := 0; i < valuesCount; i++ {
		items = append(items, "?")
	}
	pattern = strings.Join(items, ",")
	return fmt.Sprintf("(%s)", pattern)
}

func checkErr(err error) {
	if err != nil {
		telegram.Mailing(fmt.Sprintf("Программа остановилась: %s", err))
		logger.Log.Fatal(err)
		panic(err)
	}
}
