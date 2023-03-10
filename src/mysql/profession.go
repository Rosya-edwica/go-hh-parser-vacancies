package mysql

import (
	"fmt"
	"strings"

	"github.com/Rosya-edwica/go-hh-parser-vacancies/src/models"
)


func SetParsedStatusToProfession(id int) {
	db := connect()
	defer db.Close()

	query := fmt.Sprintf(`update %s set parsed=true where id=%d`, TableProfessions, id)
	fmt.Println(query)
	tx, _ := db.Begin()
	_, err := db.Exec(query)
	checkErr(err)
	tx.Commit()
} 


func GetProfessions() (professions []models.Profession) {
	db := connect()
	defer db.Close()

	query := fmt.Sprintf("SELECT id, name, other_names, level, parent_id, area_id FROM %s WHERE parsed = false", TableProfessions)
	rows, err := db.Query(query)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		var (
			name, other string
			id, level, parent_id, profRole int
		)
		err = rows.Scan(&id, &name, &other, &level, &parent_id, &profRole)
		checkErr(err)
		prof := models.Profession{
			Id: id,
			Name: name,
			OtherNames: strings.Split(other, "|"),
			Level: level,
			ParentId: parent_id,
			ProfRoleId: profRole,
		}
		professions = append(professions, prof)

	}
	return
}
