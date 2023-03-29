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

	var query string
	if CITY_LIMIT == "" {
		query = fmt.Sprintf("SELECT id_hh, id_edwica, name FROM %s WHERE id_hh != 0 ORDER BY id_hh", TableCity)
		logger.Log.Print("Ограничений на поиск по городам нет")
	} else {
		query = fmt.Sprintf("SELECT id_hh, id_edwica, name FROM %s WHERE id_hh != 0 ORDER BY id_hh LIMIT %s", TableCity, CITY_LIMIT)
		logger.Log.Printf("Ограничение на поиск по городам: %s", CITY_LIMIT)

	}
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
