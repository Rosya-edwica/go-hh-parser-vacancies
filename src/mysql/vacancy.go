package mysql

import (
	"fmt"

	"github.com/Rosya-edwica/go-hh-parser-vacancies/src/logger"
	"github.com/Rosya-edwica/go-hh-parser-vacancies/src/models"
)

func SaveOneVacancy(v models.Vacancy) (err error) {
	if v.Title == "" {
		return
	}

	db := connect()
	defer db.Close()

	columns := buildPatternInsertValues(12)
	smt := fmt.Sprintf(`INSERT INTO %s (hh_id, hh_url, name, city_id, position_id, hh_prof_areas, hh_specs, experience, salary_from, salary_to, key_skills, vacancy_date) VALUES %s`,
		TableVacancy, columns)
	tx, _ := db.Begin()
	_, err = db.Exec(smt, v.Id, v.Url, v.Title, v.CityId, v.ProfessionId, v.ProfAreas, v.Specializations, v.Experience, v.SalaryFrom, v.SalaryTo, v.Skills, v.DateUpdate)
	if err != nil {
		logger.Log.Printf("Ошибка: Вакансия %d не была добавлена в базу - %s", v.Id, err)
		err = tx.Commit()
		db.Close()
		checkErr(err)
		return
	}
	err = tx.Commit()
	checkErr(err)
	logger.Log.Printf("Успех: Вакансия %d была добавлена в базу", v.Id)
	return nil
}
