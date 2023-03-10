package mysql

import (
	"fmt"
	"os"

	"github.com/Rosya-edwica/go-hh-parser-vacancies/src/logger"
	"github.com/Rosya-edwica/go-hh-parser-vacancies/src/models"
	"github.com/joho/godotenv"
)

func GetCities() (cities []models.City) {
	db := connect()
	defer db.Close()

	err := godotenv.Load(".env")
	checkErr(err)
	CITY_LIMIT := os.Getenv("CITY_LIMIT")

	if CITY_LIMIT == "" {
		logger.Log.Fatal("Проверь переменную окружения CITY_LIMIT")
	}
	query := fmt.Sprintf("SELECT * FROM %s WHERE id_hh != 0 ORDER BY id_hh LIMIT %s", TableCity, CITY_LIMIT)
	rows, err := db.Query(query)
	checkErr(err)
	for rows.Next() {
		var name string
		var hh_id, edwica_id int
		err = rows.Scan(&hh_id, &edwica_id, &name)
		checkErr(err)
		cities = append(cities, models.City{
			HH_ID:     hh_id,
			EDWICA_ID: edwica_id,
			Name:      name,
		})
	}
	defer rows.Close()
	return
}
