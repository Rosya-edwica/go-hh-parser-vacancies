package mysql

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/Rosya-edwica/go-hh-parser-vacancies/src/models"
)

func SetParsedStatusToProfession(id int) {
	db := connect()
	defer db.Close()

	query := fmt.Sprintf(`update %s set parsed=true where position_id=%d`, TableProfessions, id)
	fmt.Println(query)
	tx, _ := db.Begin()
	_, err := db.Exec(query)
	checkErr(err)
	tx.Commit()
}

func GetProfessionsFromFile(fromFile bool) (professions []models.Profession) {
	var query string
	if fromFile {
		areas := strings.ToLower(arrayToPostgresList(readProfAreasFromFile()))
		query = `SELECT position.id, position.name, position.other_names
		FROM position
		LEFT JOIN position_to_prof_area ON position_to_prof_area.position_id=position.id
		LEFT JOIN prof_area_to_specialty ON prof_area_to_specialty.id=position_to_prof_area.area_id
		LEFT JOIN professional_area ON professional_area.id=prof_area_to_specialty.prof_area_id
		WHERE LOWER(professional_area.name) IN ` + areas
		fmt.Println("Парсим профессии этих профобластей: ", areas)
	} else {
		query = fmt.Sprintf("SELECT id, name, other_names FROM %s", TableProfessions)
		fmt.Println("Парсим абсолютно все профессии")
	}

	db := connect()
	defer db.Close()
	rows, err := db.Query(query)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		var (
			name string
			other sql.NullString
			id int
		)
		err = rows.Scan(&id, &name, &other)
		checkErr(err)

		prof := models.Profession{
			Id:         id,
			Name:       name,
			OtherNames: strings.Split(other.String, "|"),
		}
		professions = append(professions, prof)

	}
	return
}

func GetProfessionsFromAreaFile() (professions []models.Profession) {
	db := connect()
	defer db.Close()

	query := `SELECT position.name, position.id, position.other_names
	FROM position
	LEFT JOIN position_to_prof_area ON position_to_prof_area.position_id=position.id
	LEFT JOIN prof_area_to_specialty ON prof_area_to_specialty.id=position_to_prof_area.area_id
	LEFT JOIN professional_area ON professional_area.id=prof_area_to_specialty.prof_area_id
	WHERE LOWER(professional_area.name) = 'правоохранительные органы'`
	rows, err := db.Query(query)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		var (
			name string
			other sql.NullString
			id int
		)
		err = rows.Scan(&id, &name, &other)
		checkErr(err)

		prof := models.Profession{
			Id:         id,
			Name:       name,
			OtherNames: strings.Split(other.String, "|"),
		}
		professions = append(professions, prof)
	}
	return
}

func readProfAreasFromFile() (areas []string) {
	filepath := "prof_areas.txt"
	file, err := os.Open(filepath)
	if err != nil {
		panic("Создайте файл")
	}

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		areas = append(areas, fileScanner.Text())
	}
	file.Close()
	return
}

func arrayToPostgresList(items []string) (result string) {
	var updatedList []string
	for _, i := range items {
		updatedList = append(updatedList, fmt.Sprintf("'%s'", i))
	}
	result = "(" + strings.Join(updatedList, ",") + ")"
	return
}